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




// connect <i class="fa fa-database

func (p *Product) setProductImage(ctx *gin.Context, products *models.Product) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return nil
	}

	if products.Image == "" {
		//1. ตัด http://localhost:80800/uploads/prosucts/<ID>/image.png ให้เหลือ /uploads/prosucts/<ID>/image.png
		products.Image = strings.Replace(products.Image, os.Getenv("HOST"), "", 1)
		//2. แทนค่าพาธปัจจุบัน<WD>/uploads/prosucts/<ID>/image.png
		pwd, _ := os.Getwd()
		//3.remove <WD>/uploads/prosucts/<ID>/image.png
		os.Remove(pwd + products.Image)
	}

	path := "uploads/products/" + strconv.Itoa(int(products.ID))
	os.MkdirAll(path, 0755)

	filename := path + "/"+ file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	products.Image = os.Getenv("HOST") + "/" + filename

	p.DB.Save(products)

	return nil
}


func pagingResource(ctx *gin.Context, query *gorm.DB, records interface{}) *pagingResult {
	// 1. Get limit, page ?limit=10&page=2
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1")) //ดึงquery-string มาเก็บเป็น int
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "12"))
	// 2. count records
	var count int
	query.Model(records).Count(&count) //ค้นหาข้อมูลในdatabase นับด้วย Count ส่งไปที่*count

	// 3. Find Records
	// limit, offset
	// limit => 10
	// page => 1, 1 - 10, offset => 0 		// ถ้า limit มีค่า 10 , page คือ 1 ,  offset ต้องเท่ากับ 0
	// page => 2, 11 - 20, offset => 10 	// ถ้า limit มีค่า 10 , page คือ 2 ,  offset ต้องเท่ากับ 10
	// page => 3, 21 - 30, offset => 20 	// ถ้า limit มีค่า 10 , page คือ 3 ,  offset ต้องเท่ากับ 20
	offset := (page - 1) * limit
	query.Limit(limit).Offset(offset).Find(records)

	// 4. total page

	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	// 5. Find nextPage

	var nextPage int
	if nextPage == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page + 1
	}

	// 6. create pagingResult

	return &pagingResult{
		Page:      page,
		Limit:     limit,
		Count:     count,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		TotalPage: totalPage,
	}

}

      