package main

func getDefaultJsonResponse() JSONResponse {
	return JSONResponse{
		Status:               "",
		Output:               "No output",
		CompiledSuccessfully: false,
		TestPassed:           false,
	}
}
