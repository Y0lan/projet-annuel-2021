package main

func getDefaultJsonResponse() JSONResponse {
	return JSONResponse{
		Status:               "",
		CompiledSuccessfully: false,
		TestPassed:           false,
	}
}
