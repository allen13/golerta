package models

import "github.com/kataras/iris"

func StandardResponse(ctx *iris.Context, response interface{}, err error){
	if err != nil {
		ctx.JSON(iris.StatusInternalServerError, ErrorResponse{Status: "error", Message: err.Error()})
	}
	ctx.JSON(iris.StatusOK, response)
}
