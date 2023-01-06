package main

// func (app *application) InvalidCredentials(w http.ResponseWriter, r *http.Request, password, repeatPassword string) error {
// 	if password != repeatPassword {
// 		var msg string = ""
// 		data := make(map[string]interface{})
// 		data["message"] = msg
// 		if err := app.renderTemplate(w, r, "error", data); err != nil {
// 			app.errorLog.Println(err)
// 			return err
// 		}
// 	}
// }
//  REWRITE INVALID CREDENTIALS, TO CHECK HASH AND PASSWORD
//  CREATE FUNCTION, WHICH CHECK IF USER OR EMAIL ALREADY IS IN DATABASE
