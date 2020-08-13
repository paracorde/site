package utils

import (
   "path/filepath"

   "github.com/flosch/pongo2"
   _ "github.com/flosch/pongo2-addons"
)

const (
   templateDir = "templates"
)

func LoadTemplates() (map[string]*pongo2.Template, error) {
   templates := make(map[string]*pongo2.Template)
   var err error
   templates["base"], err = pongo2.FromFile(filepath.Join(templateDir, "base.html.pt"))
   templates["post"], err = pongo2.FromFile(filepath.Join(templateDir, "post.html.pt"))
   templates["modifypost"], err = pongo2.FromFile(filepath.Join(templateDir, "modifypost.html.pt"))
   templates["displayposts"], err = pongo2.FromFile(filepath.Join(templateDir, "displayposts.html.pt"))
   templates["makeannouncement"], err = pongo2.FromFile(filepath.Join(templateDir, "makeannouncement.html.pt"))
   templates["calendar"], err = pongo2.FromFile(filepath.Join(templateDir, "calendar.html.pt"))
   if err != nil {
      return nil, err
   }
   return templates, nil
}

func getPostFromKey(in *pongo2.Value, key *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
   m := in.Interface().(map[string]*Post)
   return pongo2.AsValue(m[key.String()]), nil
}

func getAnnouncementFromKey(in *pongo2.Value, key *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
   m := in.Interface().(map[string]*Announcement)
   return pongo2.AsValue(m[key.String()]), nil
}

func init() {
   pongo2.RegisterFilter("postkey", getPostFromKey)
   pongo2.RegisterFilter("announcementkey", getAnnouncementFromKey)
}
