go mod init kp-app

   gin = https://github.com/gin-gonic/gin    


.vscode 
{
    "go.useLanguageServer": true,
    "[go]": {
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": true,
        },
        // Optional: Disable snippets, as they conflict with completion ranking.
        "editor.snippetSuggestions": "none",
    },
    "[go.mod]": {
        "editor.formatOnSave": true,
        "editor.codeActionsOnSave": {
            "source.organizeImports": true,
        },
    },
    "gopls": {
        // Add parameter placeholders when completing a function.
        "usePlaceholders": true,

        // If true, enable additional analyses with staticcheck.
        // Warning: This will significantly increase memory usage.
        "staticcheck": false,
    }
}

//upload file
// Get file
		file, _ := ctx.FormFile("image")

		// Create file
		path := "uploads/products/" + strconv.Itoa(int(p.ID)) // ID => 8, uploads/articles/8/image.png
		os.MkdirAll(path, 0755)                               // -> uploads/products/8

		// Upload file
		filename := path + file.Filename
		if err := ctx.SaveUploadedFile(file, filename); err != nil {
			log.Fatal(err.Error())
		}

		// Attach file to products
		p.Image = "http://localhost:8080/" + filename
      