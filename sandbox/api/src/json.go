package main

func getDefaultJsonResponse() JSONResponse {
	return JSONResponse{
		Status:               "",
		Output:               "No output",
		CompiledSuccessfully: false,
		TestPassed:           false,
		CodeQuality:          80.0,
	}
}

func populateJsonResponse(result CodeData, response JSONResponse) JSONResponse {
	var testPassed = result.ExitSuccessfully && response.CompiledSuccessfully
	var testDidNotPass = !result.ExitSuccessfully && response.CompiledSuccessfully
	var compilationPassed = result.ExitSuccessfully && !response.CompiledSuccessfully
	var compilationDidNotPass = !result.ExitSuccessfully && !response.CompiledSuccessfully

	if testPassed {
		response.Status = "success"
		response.Output = response.Output + result.CommandOutput
		response.TestPassed = true
	}

	if compilationPassed {
		response.Status = "success"
		response.CompiledSuccessfully = true
		response.Output = result.CommandOutput
		response.Duration = result.Duration
		response.CodeQuality = result.CodeQuality
	}

	if testDidNotPass {
		response.Status = "failure"
		response.TestPassed = false
		response.Error = result.CommandOutput
		response.CodeQuality = 0
	}

	if compilationDidNotPass {
		response.Status = "failure"
		response.Error = result.CommandOutput
		response.CompiledSuccessfully = false
		response.Output = ""
		response.Error = result.CommandOutput
		response.Duration = result.Duration
	}
	return response
}
