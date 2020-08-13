// TO-DO
// figure out cleaner way to include the flashed message in every pongo context

package main

import (
   // "fmt"
   "io/ioutil"
   "log"
   "net/http"
   // "path/filepath"

   "github.com/paracorde/ssite/utils"

   "golang.org/x/crypto/bcrypt"
   "github.com/labstack/echo/v4"
   "github.com/flosch/pongo2"
   _ "github.com/flosch/pongo2-addons"
)

const (
   templateDir = "templates"
)

// TemplateSet holding all templates
// var templateSet = pongo2.NewSet("templates")
var templates = make(map[string]*pongo2.Template)

// Objects to hold all posts/announcements
var posts *utils.Posts
var announcements *utils.Announcements

// Password for editing/deleting/creating
var hashed_password []byte

// e.GET("/post/:id", viewPost)
func viewPost(context echo.Context) error {
   postID := context.Param("id")
   post, present := posts.PostMap[postID]
   if !present {
      // redirect to 404 not found
      return nil
   }
   html, err := templates["post"].Execute(pongo2.Context{
      "PostID": post.PostID.String(),
      "Title": post.Title,
      "Author": post.Author,
      "Timestamp": post.Timestamp,
      "Body": post.Body,
   })
   if err != nil {
      return err
   }
   return context.HTML(http.StatusOK, html)
}

// e.GET("/create-post", createPost)
func createPost(context echo.Context) error {
   html, err := templates["modifypost"].Execute(pongo2.Context{
      "FormAction": "/create-post",
      "Title": "",
      "Author": "",
      "Body": "",
      "FormPurpose": "Submit post",
   })
   if err != nil {
      return err
   }
   return context.HTML(http.StatusOK, html)
}

// e.POST("/create-post", submitPost)
func submitPost(context echo.Context) error {
   title := context.FormValue("title")
   author := context.FormValue("author")
   body := context.FormValue("body")
   password := []byte(context.FormValue("password"))
   err := bcrypt.CompareHashAndPassword(hashed_password, password)
   if err != nil {
      utils.Flash(context, "Password entered is incorrect.")
      html, err := templates["modifypost"].Execute(pongo2.Context{
         "FormAction": "/create-post",
         "Title": title,
         "Author": author,
         "Body": body,
         "FormPurpose": "Submit post",
      })
      if err != nil {
         return err
      }
      return context.HTML(http.StatusOK, html)
   }
   post, err := utils.PublishPost(title, author, body)
   if err != nil {
      return err
   }
   posts.AddPost(post)
   utils.Flash(context, "Your post has successfully been submitted!")
   return context.Redirect(http.StatusFound, "/posts/"+post.PostID.String())
}

// e.GET("/edit/:id", editPost)
func editPost(context echo.Context) error {
   postID := context.Param("id")
   post, present := posts.PostMap[postID]
   if !present {
      // redirect to 404 not found
      return nil
   }
   html, err := templates["modifypost"].Execute(pongo2.Context{
      "FormAction": "/edit/"+postID,
      "Title": post.Title,
      "Author": post.Author,
      "Body": post.Body,
      "FormPurpose": "Commit changes",
   })
   if err != nil {
      return err
   }
   return context.HTML(http.StatusOK, html)
}

// e.POST("/edit/:id", commitChanges)
func commitChanges(context echo.Context) error {
   postID := context.Param("id")
   post, present := posts.PostMap[postID]
   if !present {
      // redirect to 404 not found
      return nil
   }
   title := context.FormValue("title")
   author := context.FormValue("author")
   body := context.FormValue("body")
   password := []byte(context.FormValue("password"))
   err := bcrypt.CompareHashAndPassword(hashed_password, password)
   if err != nil {
      utils.Flash(context, "Password entered is incorrect.")
      html, err := templates["modifypost"].Execute(pongo2.Context{
         "FormAction": "/edit/"+postID,
         "Title": title,
         "Author": author,
         "Body": body,
         "FormPurpose": "Commit changes",
      })
      if err != nil {
         return err
      }
      return context.HTML(http.StatusOK, html)
   }
   post.Title = title
   post.Author = author
   post.Body = body
   err = post.WritePost()
   if err != nil {
      return err
   }
   utils.Flash(context, "Post successfully edited!")
   return context.Redirect(http.StatusFound, "/posts/"+post.PostID.String())
}

// e.GET("/display-posts", displayPosts)
func displayPosts(context echo.Context) error {
   html, err := templates["displayposts"].Execute(pongo2.Context{
      "PostSlice": posts.PostSlice,
      "PostMap": posts.PostMap,
      "WebM": false,
   })
   if err != nil {
      return err
   }
   return context.HTML(http.StatusOK, html)
}

// e.GET("/make-announcement", makeAnnouncement)
func makeAnnouncement(context echo.Context) error {
   html, err := templates["makeannouncement"].Execute(pongo2.Context{
      "Body": "",
   })
   if err != nil {
      return err
   }
   return context.HTML(http.StatusOK, html)
}

// e.POST("/make-announcement", pushAnnouncement)
func pushAnnouncement(context echo.Context) error {
   body := context.FormValue("body")
   password := []byte(context.FormValue("password"))
   err := bcrypt.CompareHashAndPassword(hashed_password, password)
   if err != nil {
      utils.Flash(context, "Password entered is incorrect.")
      html, err := templates["makeannouncement"].Execute(pongo2.Context{
         "Body": body,
      })
      if err != nil {
         return err
      }
      return context.HTML(http.StatusOK, html)
   }
   announcement, err := utils.PublishAnnouncement(body)
   if err != nil {
      return err
   }
   announcements.AddAnnouncement(announcement)
   utils.Flash(context, "Announcement successfully made!")
   return context.Redirect(http.StatusFound, "/calendar")
}

// e.GET("/calendar", showCalendar)
func showCalendar(context echo.Context) error {
   slice := announcements.AnnouncementSlice
   html, err := templates["calendar"].Execute(pongo2.Context{
      "AnnouncementSlice": slice[utils.Max(len(slice)-10, 0):],
      "AnnouncementMap": announcements.AnnouncementMap,
   })
   if (err != nil) {
      return err
   }
   return context.HTML(http.StatusOK, html)
}

// e.GET("/webmaster-posts", webmasterDisplayPosts)
func webmasterDisplayPosts(context echo.Context) error {
   html, err := templates["displayposts"].Execute(pongo2.Context{
      "PostSlice": posts.PostSlice,
      "PostMap": posts.PostMap,
      "WebM": true,
   })
   if err != nil {
      return err
   }
   return context.HTML(http.StatusOK, html)
}

func main() {
   var err error
   posts, err = utils.LoadPosts()
   if err != nil {
      log.Fatal(err)
   }
   announcements, err = utils.LoadAnnouncements()
   if err != nil {
      log.Fatal(err)
   }
   templates, err = utils.LoadTemplates()
   if err != nil {
      log.Fatal(err)
   }
   hashed_password, err = ioutil.ReadFile("password")
   if err != nil {
      log.Fatal(err)
   }
   e := echo.New()
   e.Static("/static", "static")
   e.GET("/posts/:id", viewPost)
   e.GET("/create-post", createPost)
   e.POST("/create-post", submitPost)
   e.GET("/edit/:id", editPost)
   e.POST("/edit/:id", commitChanges)
   e.GET("/display-posts", displayPosts)
   e.GET("/make-announcement", makeAnnouncement)
   e.POST("/make-announcement", pushAnnouncement)
   e.GET("/calendar", showCalendar)
   e.Logger.Fatal(e.Start(":8080"))
}
