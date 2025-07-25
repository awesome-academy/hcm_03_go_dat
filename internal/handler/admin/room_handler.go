package admin

import (
	"fmt"
	"hotel-management/internal/constant"
	"hotel-management/internal/dto"
	"hotel-management/internal/usecase/admin_usecase"
	"hotel-management/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomUseCase *admin_usecase.RoomUseCase
}

func NewRoomHandler(roomUseCase *admin_usecase.RoomUseCase) *RoomHandler {
	return &RoomHandler{roomUseCase: roomUseCase}
}
func (h *RoomHandler) RoomManagementPage(c *gin.Context) {
	rooms, err := h.roomUseCase.GetAllRooms(c.Request.Context())
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": utils.T(c, "error.failed_to_load_rooms")})
		return
	}
	c.HTML(http.StatusOK, "room.html", gin.H{
		"Title": "Room management",
		"Rooms": rooms,
	})
}

func (h *RoomHandler) CreateRoomPage(c *gin.Context) {
	c.HTML(http.StatusOK, "create_room.html", gin.H{
		"Title": "Create room",
	})
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	roomType := strings.TrimSpace(c.PostForm("type"))
	priceStr := c.PostForm("price_per_night")
	bedStr := c.PostForm("bed_num")
	viewType := strings.TrimSpace(c.PostForm("view_type"))
	description := strings.TrimSpace(c.PostForm("description"))
	hasAircon := c.PostForm("has_aircon") == "on"
	isAvailable := c.PostForm("is_available") == "on"

	// Validate
	if name == "" || roomType == "" || priceStr == "" || bedStr == "" || viewType == "" {
		c.HTML(http.StatusBadRequest, "create_room.html", gin.H{"error": utils.T(c, "error.invalid_request")})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price < 0 {
		c.HTML(http.StatusBadRequest, "create_room.html", gin.H{"error": utils.T(c, "error.invalid_price_per_night")})
		return
	}

	beds, err := strconv.Atoi(bedStr)
	if err != nil || beds < 1 {
		c.HTML(http.StatusBadRequest, "create_room.html", gin.H{"error": utils.T(c, "error.invalid_bed_num")})
		return
	}

	// Parse files
	form, err := c.MultipartForm()
	if err != nil {
		c.HTML(http.StatusBadRequest, "create_room.html", gin.H{"error": utils.T(c, "error.invalid_request")})
		return
	}
	files := form.File["images"]

	createRoomRequest := &dto.CreateRoomRequest{
		Name:          name,
		Type:          roomType,
		PricePerNight: price,
		BedNum:        beds,
		HasAircon:     hasAircon,
		ViewType:      viewType,
		Description:   description,
		IsAvailable:   isAvailable,
		ImageFiles:    files,
	}

	err = h.roomUseCase.CreateRoom(c, createRoomRequest)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "create_room.html", gin.H{"error": utils.T(c, err.Error())})
		return
	}

	c.Redirect(http.StatusSeeOther, constant.RoomManagementPath)
}

func (h *RoomHandler) RoomDetailPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": utils.T(c, "error.invalid_room_id")})
		return
	}
	room, err := h.roomUseCase.GetRoomByID(c.Request.Context(), id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": utils.T(c, "error.room_not_found")})
		return
	}

	c.HTML(http.StatusOK, "room_detail.html", gin.H{
		"Title": "Room detail",
		"Room":  room,
	})
}

func (h *RoomHandler) EditRoomPage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": utils.T(c, "error.invalid_room_id")})
		return
	}
	room, err := h.roomUseCase.GetRoomByID(c.Request.Context(), id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": utils.T(c, "error.room_not_found")})
		return
	}

	c.HTML(http.StatusOK, "edit_room.html", gin.H{
		"Title": "Edit room",
		"Room":  room,
	})
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "edit_room.html", gin.H{"error": utils.T(c, "error.invalid_room_id")})
		return
	}

	name := strings.TrimSpace(c.PostForm("name"))
	roomType := strings.TrimSpace(c.PostForm("type"))
	priceStr := c.PostForm("price_per_night")
	bedStr := c.PostForm("bed_num")
	viewType := strings.TrimSpace(c.PostForm("view_type"))
	description := strings.TrimSpace(c.PostForm("description"))
	hasAircon := false
	if c.PostForm("has_aircon") == "on" {
		hasAircon = true
	}
	isAvailable := false
	if c.PostForm("is_available") == "on" {
		isAvailable = true
	}

	// Validate
	if name == "" || roomType == "" || priceStr == "" || bedStr == "" || viewType == "" {
		c.HTML(http.StatusBadRequest, "edit_room.html", gin.H{"error": utils.T(c, "error.invalid_request")})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price < 0 {
		c.HTML(http.StatusBadRequest, "edit_room.html", gin.H{"error": utils.T(c, "error.invalid_price_per_night")})
		return
	}

	beds, err := strconv.Atoi(bedStr)
	if err != nil || beds < 1 {
		c.HTML(http.StatusBadRequest, "edit_room.html", gin.H{"error": utils.T(c, "error.invalid_bed_num")})
		return
	}

	// Parse Multipart
	form, err := c.MultipartForm()
	if err != nil {
		c.HTML(http.StatusBadRequest, "edit_room.html", gin.H{"error": utils.T(c, "error.invalid_request")})
		return
	}

	files := form.File["images"]

	// Parse deleted image id list
	deletedIDs := c.PostFormArray("delete_image_ids")

	var deletedImageIDs []int
	for _, idStr := range deletedIDs {
		id, err := strconv.Atoi(idStr)
		if err == nil {
			deletedImageIDs = append(deletedImageIDs, id)
		}
	}

	updateReq := &dto.EditRoomRequest{
		ID:            roomID,
		Name:          name,
		Type:          roomType,
		PricePerNight: price,
		BedNum:        beds,
		HasAircon:     hasAircon,
		ViewType:      viewType,
		Description:   description,
		IsAvailable:   isAvailable,
		ImageFiles:    files,
		ImageDeletes:  deletedImageIDs,
	}
	fmt.Println("updateReq:", updateReq)

	err = h.roomUseCase.UpdateRoom(c, updateReq)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "edit_room.html", gin.H{"error": utils.T(c, err.Error())})
		return
	}

	c.Redirect(http.StatusSeeOther, constant.RoomManagementPath)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": utils.T(c, "error.invalid_room_id")})
		return
	}
	err = h.roomUseCase.DeleteRoom(c.Request.Context(), id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"error": utils.T(c, err.Error())})
		return
	}

	c.Redirect(http.StatusSeeOther, constant.RoomManagementPath)
}
