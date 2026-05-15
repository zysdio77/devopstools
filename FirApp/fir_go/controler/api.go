package controler

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

var allowedExtensions = map[string]bool{}

func initAllowedExtensions() {
	if len(allowedExtensions) > 0 {
		return
	}
	exts := g.Cfg().MustGet(gctx.New(), "server.allowedExtensions").Strings()
	for _, ext := range exts {
		allowedExtensions[strings.ToLower(ext)] = true
	}
	if len(allowedExtensions) == 0 {
		allowedExtensions[".ipa"] = true
		allowedExtensions[".apk"] = true
	}
}

func validateFilename(name string) error {
	clean := filepath.Clean(name)
	if clean != name && clean != filepath.Base(name) {
		return fmt.Errorf("invalid filename: path traversal detected")
	}
	if strings.Contains(name, "..") {
		return fmt.Errorf("invalid filename: path traversal detected")
	}
	return nil
}

func validateFileExtension(filename string) error {
	initAllowedExtensions()
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("file type %s is not allowed", ext)
	}
	return nil
}

func safeJoinPath(basePath, elem string) (string, error) {
	if err := validateFilename(elem); err != nil {
		return "", err
	}
	fullPath := filepath.Join(basePath, elem)
	cleanBase := filepath.Clean(basePath)
	if !strings.HasPrefix(filepath.Clean(fullPath), cleanBase) {
		return "", fmt.Errorf("path traversal detected")
	}
	return fullPath, nil
}

func deleteFileSafe(basePath, filename string) error {
	path, err := safeJoinPath(basePath, filename)
	if err != nil {
		return err
	}
	if gfile.IsFile(path) {
		return gfile.Remove(path)
	}
	return nil
}

func GetInfo(r *ghttp.Request) {
	ctx := r.Context()
	name := r.Get("name")
	result, err := SelectDb(ctx, name)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	g.Log().Info(ctx, "method:GetInfo", result)

	var data interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		r.Response.WriteJson(g.Map{"success": false, "message": "data error"})
		return
	}
	r.Response.WriteJson(data)
}

func DeleteInfo(r *ghttp.Request) {
	ctx := r.Context()
	name := r.Get("name").String()
	systemtype := r.Get("system_type").String()

	if err := validateFilename(name); err != nil {
		r.Response.WriteJson(g.Map{"success": false, "message": err.Error()})
		return
	}
	if err := validateFilename(systemtype); err != nil {
		r.Response.WriteJson(g.Map{"success": false, "message": err.Error()})
		return
	}

	result, err := DeleteDb(ctx, name)
	if err != nil {
		g.Log().Error(ctx, "method:DeleteInfo", err)
		r.Response.WriteJson(g.Map{"success": false, "message": err.Error()})
		return
	}

	templateRoot := g.Cfg().MustGet(ctx, "server.template").String()
	filePath := fmt.Sprintf("%s/%s", systemtype, name)
	deleteFileSafe(templateRoot, filePath)

	if systemtype == "ios" {
		plistName := strings.Replace(name, ".ipa", ".plist", 1)
		plistPath := fmt.Sprintf("plist/%s", plistName)
		deleteFileSafe(templateRoot, plistPath)
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"lines":   result,
	})
}

func IosPlist(ctx context.Context, filename string) error {
	ipauri, err := g.Cfg().Get(ctx, "server.ipauri")
	if err != nil {
		return err
	}
	templateDir, err := g.Cfg().Get(ctx, "server.template")
	if err != nil {
		return err
	}
	templatename, err := g.Cfg().Get(ctx, "server.templatename")
	if err != nil {
		return err
	}

	plistName := strings.Replace(filename, ".ipa", ".plist", 1)
	plistPath := fmt.Sprintf("%v/plist/%v", templateDir, plistName)
	templatePath := fmt.Sprintf("%v/%v", templateDir, templatename)

	if err := gfile.Copy(templatePath, plistPath); err != nil {
		return err
	}

	replaceContent := fmt.Sprintf("%v%v", ipauri, filename)
	if err := gfile.ReplaceFile("ipaurl", replaceContent, plistPath); err != nil {
		return err
	}
	return nil
}

func Upload(r *ghttp.Request) {
	ctx := r.Context()
	templateRoot := g.Cfg().MustGet(ctx, "server.template").String()
	uploadType := r.Get("system_type").String()
	note := r.Get("note")

	if err := validateFilename(uploadType); err != nil {
		r.Response.WriteJson(g.Map{"success": false, "message": err.Error()})
		return
	}

	saveAddr := filepath.Join(templateRoot, uploadType)
	files := r.GetUploadFiles("upload-file")
	if len(files) == 0 {
		r.Response.WriteJson(g.Map{"success": false, "message": "no file uploaded"})
		return
	}

	uploadedFile := files[0]
	if err := validateFileExtension(uploadedFile.Filename); err != nil {
		r.Response.WriteJson(g.Map{"success": false, "message": err.Error()})
		return
	}
	if err := validateFilename(uploadedFile.Filename); err != nil {
		r.Response.WriteJson(g.Map{"success": false, "message": err.Error()})
		return
	}

	names, err := files.Save(saveAddr)
	if err != nil {
		g.Log().Error(ctx, "upload failed:", err)
		r.Response.WriteJson(g.Map{"success": false, "message": err.Error()})
		return
	}

	g.Log().Info(ctx, "upload successfully:", names, "saveAddr:", saveAddr)

	info := &Infos{
		Name:       names[0],
		Systemtype: r.Get("system_type").String(),
		Note:       note.String(),
	}

	result, err := InsertDb(ctx, info)
	if err != nil {
		g.Log().Error(ctx, "InsertDb failed:", err)
		r.Response.WriteJson(g.Map{"success": false, "message": err.Error()})
		return
	}

	if uploadType == "ios" {
		if err := IosPlist(ctx, names[0]); err != nil {
			g.Log().Error(ctx, "IosPlist failed:", err)
		}
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"lines":   result,
		"message": info,
	})
}

func GetPage(r *ghttp.Request) {
	ctx := r.Context()
	name := r.Get("name")
	page := r.Get("page")

	total, pages, resultStr, err := PagesData(ctx, name, page)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var records []map[string]interface{}
	if err := json.Unmarshal([]byte(resultStr), &records); err != nil {
		g.Log().Error(ctx, "JSON parse failed:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "data parse error",
		})
		return
	}

	g.Log().Info(ctx, "method:GetPage", g.Map{"page": pages, "total": total})

	r.Response.WriteJson(g.Map{
		"total":  total,
		"page":   pages,
		"result": records,
	})
}

func UpdateInfo(r *ghttp.Request) {
	ctx := r.Context()
	name := r.Get("name").String()
	systemtype := r.Get("system_type").String()

	result, err := UpdateDb(ctx, name, systemtype)
	if err != nil {
		g.Log().Error(ctx, "method:UpdateInfo", err)
		r.Response.WriteJson(g.Map{"success": false, "message": err.Error()})
		return
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"lines":   result,
	})
}
