package handler

import (
	"strconv"

	"github.com/fahri/go-tije/internal/service"
	"github.com/gofiber/fiber/v2"
)

type VehicleHandler struct {
	service service.VehicleService
}

func NewVehicleHandler(service service.VehicleService) *VehicleHandler {
	return &VehicleHandler{
		service: service,
	}
}

func (h *VehicleHandler) GetLatestLocation(c *fiber.Ctx) error {
	vehicleID := c.Params("vehicle_id")
	if vehicleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "vehicle_id is required",
		})
	}
	
	location, err := h.service.GetLatestLocation(c.Context(), vehicleID)
	if err != nil {
		if err.Error() == "vehicle not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get location",
		})
	}
	
	return c.JSON(location)
}

func (h *VehicleHandler) GetLocationHistory(c *fiber.Ctx) error {
	vehicleID := c.Params("vehicle_id")
	if vehicleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "vehicle_id is required",
		})
	}
	
	startStr := c.Query("start")
	endStr := c.Query("end")
	
	if startStr == "" || endStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "start and end timestamps are required",
		})
	}
	
	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid start timestamp",
		})
	}
	
	end, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid end timestamp",
		})
	}
	
	locations, err := h.service.GetLocationHistory(c.Context(), vehicleID, start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get location history",
		})
	}
	
	return c.JSON(locations)
}