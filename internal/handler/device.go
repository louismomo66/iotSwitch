package handler

import (
	"encoding/json"
	"errors"
	"iot_switch/internal/models"
	"iot_switch/internal/repository"
	"iot_switch/internal/utils"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type DeviceController struct {
	Repo repository.DeviceRepository
}

func NewDeviceController(repo repository.DeviceRepository) *DeviceController {
	return &DeviceController{Repo: repo}
}

func (d *DeviceController) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	// log.Println("Trying to register")
	// requestBody, _ := ioutil.ReadAll(r.Body)
    //     defer r.Body.Close()
    //     log.Printf("Failed to decode. Request body: %s\n", requestBody)

    var device models.Device
    if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		
        utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    existingDevice, err := d.Repo.GetDeviceByESP32ID(device.ESP32ID)
    if err != nil {
        if err.Error() == "record not found" {
            if err := d.Repo.CreateDevice(&device); err != nil {
				log.Println("Failed to register")
                utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to register ESP32")
                return
            }
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(device)
            return
        } else {
			log.Println("Table error")
            utils.WriteJSONError(w, http.StatusInternalServerError, err, "Device table error")
            return
        }
    }

    existingDevice.NumRelays = device.NumRelays

    for _, relay := range device.Relays {
        existingRelay, err := d.Repo.GetRelayByESP32IDAndPin(device.ESP32ID, relay.Pin)
        if err != nil {
            if err.Error() == "record not found" {
                existingDevice.Relays = append(existingDevice.Relays, relay)
            } else {
                utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to fetch existing relays")
                return
            }
        } else {
            existingRelay.State = relay.State
            if err := d.Repo.UpdateRelayState(existingRelay); err != nil {
                utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to update relay")
                return
            }
        }
    }

    if err := d.Repo.UpdateDevice(existingDevice); err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to update ESP32")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(existingDevice)
}

func (d *DeviceController) SetRelayState(w http.ResponseWriter, r *http.Request) {
    var req struct {
        ESP32ID string `json:"esp32_id"`
        Pin     int    `json:"pin"`
        State   bool   `json:"state"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.WriteJSONError(w, http.StatusBadRequest, err, "Invalid request payload")
        return
    }

    device, err := d.Repo.GetDeviceByESP32ID(req.ESP32ID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            utils.WriteJSONError(w, http.StatusNotFound, err, "Device not found")
        } else {
            utils.WriteJSONError(w, http.StatusInternalServerError, err, "Error fetching device")
        }
        return
    }

    var relay *models.Relay
    for _, r := range device.Relays {
        if r.Pin == req.Pin {
            relay = &r
            break
        }
    }

    if relay == nil {
        utils.WriteJSONError(w, http.StatusNotFound, errors.New("relay not found"), "relay not found")
        return
    }

    relay.State = "off"
    if req.State {
        relay.State = "on"
    }

    if err := d.Repo.UpdateRelayState(relay); err != nil {
        utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to update relay state")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(relay)
}



type RelayStates struct {
	Relays map[int]bool `json:"relays"`
}

func (h *ScheduleHandler) GetRelayStates(w http.ResponseWriter, r *http.Request) {
	esp32ID := mux.Vars(r)["esp32_id"]

	var device models.Device
	if err := h.DB.Preload("Relays").Where("esp32_id = ?", esp32ID).First(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.WriteJSONError(w, http.StatusNotFound, err, "Device not found")
		} else {
			utils.WriteJSONError(w, http.StatusInternalServerError, err, "Error fetching device")
		}
		return
	}

	var schedules []models.Schedule
	if err := h.DB.Where("relay_id IN ?", getRelayIDs(device.Relays)).Find(&schedules).Error; err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Error fetching schedules")
		return
	}

	relayStates := make(map[int]bool)
	now := time.Now()

	// Initialize all relay states to 'off'
    for _, relay := range device.Relays {
        relayStates[relay.Pin] = false
    }

    // Update states based on the current state from the database
    for _, relay := range device.Relays {
        if relay.State == "on" {
            relayStates[relay.Pin] = true
        }
    }

	// When checking schedules
for _, schedule := range schedules {
    if schedule.Active {
        start := schedule.StartTime
        end := start.Add(time.Duration(schedule.Duration) * time.Second)
        log.Printf("Checking schedule: %v, start: %v, end: %v, now: %v", schedule.ID, start, end, now)
        if now.After(start) && now.Before(end) {
            for _, relay := range device.Relays {
                if relay.ID == schedule.RelayID {
                    relayStates[relay.Pin] = true
                    log.Printf("Schedule %v activates relay pin %v", schedule.ID, relay.Pin)
                    break
                }
            }
        }
    }
}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(relayStates)
}


func getRelayIDs(relays []models.Relay) []int {
	var ids []int
	for _, relay := range relays {
		ids = append(ids, relay.ID)
	}
	return ids
}

func (d *DeviceController) GetAllDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := d.Repo.GetAllDevices()
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err, "Failed to fetch devices")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}




func (d *DeviceController) GetRelaysByESP32ID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    esp32ID := vars["esp32_id"]

    device, err := d.Repo.GetRelayByESP32ID(esp32ID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            utils.WriteJSONError(w, http.StatusNotFound, err, "Device not found")
        } else {
            utils.WriteJSONError(w, http.StatusInternalServerError, err, "Error fetching device")
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(device.Relays)
}

func (d *DeviceController) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}  
	esp32ID := mux.Vars(r)["esp32_id"]

    if err := d.Repo.DeleteDeviceByESP32ID(esp32ID); err != nil {
        if err == gorm.ErrRecordNotFound {
            utils.WriteJSONError(w, http.StatusNotFound, err, "Device not found")
        } else {
            utils.WriteJSONError(w, http.StatusInternalServerError, err, "Error deleting device")
        }
        return
    }else{
		log.Printf("Device with ID: %s deleted successfully", esp32ID)
	}
	response := map[string]string{
        "message": "Device deleted successfully",
    }
	utils.WriteJSON(w, http.StatusOK, response)

    // w.WriteHeader(http.StatusNoContent)
}
