package handler

import (
	"log"
	"net/http"
	"pocketbase-seal/helper"
	"pocketbase-seal/model/genericresponse"
	"pocketbase-seal/model/seal"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type sealHandler struct {
	db *pocketbase.PocketBase
}

func NewSealHandler(db *pocketbase.PocketBase) *sealHandler {
	return &sealHandler{
		db,
	}
}

func (h *sealHandler) GetSpecificSeal(e *core.RequestEvent) error {
	id := e.Request.PathValue("id")

	seals := seal.Seal{}

	err := h.db.DB().NewQuery("SELECT id, name, color, gender, weight, age, dob FROM seals WHERE id = {:id}").Bind(dbx.Params{
		"id": id,
	}).One(&seals)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			errResponse := genericresponse.GenericErrorResponse{Error: err, ResponseMessage: "No Record Found", ResponseCode: 404}
			return e.JSON(http.StatusNotFound, errResponse)
		}
		errResponse := genericresponse.GenericErrorResponse{Error: err, ResponseMessage: "Internal Server Error", ResponseCode: 500}
		return e.JSON(http.StatusInternalServerError, errResponse)
	}

	response := genericresponse.GenericResponse{Data: seals, ResponseMessage: "Success", ResponseCode: 200}

	return e.JSON(http.StatusOK, response)
}

func (h *sealHandler) GetSeals(e *core.RequestEvent) error {
	seals := []seal.Seal{}

	page, limit := helper.GetPagesLimit(e)

	offset := (page - 1) * limit

	err := h.db.DB().NewQuery("SELECT id, name, color, gender, weight, age, dob FROM seals ORDER BY id LIMIT {:limit} OFFSET {:offset}").Bind(dbx.Params{
		"limit":  limit,
		"offset": offset,
	}).All(&seals)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			errResponse := genericresponse.GenericErrorResponse{Error: err, ResponseMessage: "No Record Found", ResponseCode: 404}
			return e.JSON(http.StatusNotFound, errResponse)
		}
		errResponse := genericresponse.GenericErrorResponse{Error: err, ResponseMessage: "Internal Server Error", ResponseCode: 500}
		return e.JSON(http.StatusInternalServerError, errResponse)
	}

	response := genericresponse.GenericResponse{Data: seals, ResponseMessage: "Success", ResponseCode: 200}

	return e.JSON(http.StatusOK, response)

}

func (h *sealHandler) UpdateSeal(e *core.RequestEvent) error {

	seal := seal.Seal{}

	if err := e.BindBody(&seal); err != nil {
		errResponse := genericresponse.GenericErrorResponse{Error: err, ResponseMessage: "Failed reading request body", ResponseCode: 400}
		return e.JSON(http.StatusBadRequest, errResponse)
	}

	_, err := h.db.DB().NewQuery("UPDATE seals SET name = {:name}, color = {:color}, gender = {:gender}, weight = {:weight}, age = {:age}, dob = {:dob} WHERE id = {:id}").Bind(dbx.Params{
		"name":   seal.Name,
		"color":  seal.Color,
		"gender": seal.Gender,
		"weight": seal.Weight,
		"age":    seal.Age,
		"dob":    seal.Dob,
		"id":     seal.Id,
	}).Execute()

	if err != nil {
		log.Print(err)
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Failed updating seal", ResponseCode: 500}
		return e.JSON(http.StatusInternalServerError, errResponse)
	}

	response := genericresponse.GenericResponse{Data: seal, ResponseMessage: "Success", ResponseCode: 200}

	return e.JSON(http.StatusOK, response)
}

func (h *sealHandler) CreateNewSeal(e *core.RequestEvent) error {
	newSeal := seal.SealAdd{}

	if err := e.BindBody(&newSeal); err != nil {
		errResponse := genericresponse.GenericErrorResponse{Error: err, ResponseMessage: "Failed reading request body", ResponseCode: 400}
		return e.JSON(http.StatusBadRequest, errResponse)
	}

	empty := helper.CheckEmptyStruct(newSeal)

	if empty {
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Body contains empty field", ResponseCode: 400}
		return e.JSON(http.StatusBadRequest, errResponse)
	}

	_, err := h.db.DB().NewQuery("INSERT INTO seals(name, color, gender, weight, age, dob) VALUES({:name}, {:color}, {:gender}, {:weight}, {:age}, {:dob})").Bind(dbx.Params{
		"name":   newSeal.Name,
		"color":  newSeal.Color,
		"gender": newSeal.Gender,
		"weight": newSeal.Weight,
		"age":    newSeal.Age,
		"dob":    newSeal.Dob,
	}).Execute()

	if err != nil {
		log.Print(err)
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Failed creating new seal", ResponseCode: 500}
		return e.JSON(http.StatusInternalServerError, errResponse)
	}

	response := genericresponse.GenericResponse{Data: newSeal, ResponseMessage: "Success", ResponseCode: 200}
	return e.JSON(http.StatusOK, response)

}
