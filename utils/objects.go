package utils

import (
   "encoding/json"
   // "fmt"
   "io/ioutil"
   "os"
   "path/filepath"
   "strings"
   // "time"

   "github.com/segmentio/ksuid"
)

const (
   timestampFormat = "Monday, 2 January 2006 – 3:04 pm"
   timestampShortFormat = "January 2006"
   docDir = "posts"
   announcementDir = "announcements"
)

type Post struct {
   PostID         ksuid.KSUID `json:"-"` // don't include this field when marshalling
   Title          string
   Timestamp      string `json:"-"`
   TimestampShort string `json:"-"`
   Author         string
   Body           string
}

func PublishPost(title, author, body string) (*Post, error) {
   postID := ksuid.New()
   post := &Post{
      PostID: postID,
      Title: title,
      Author: author,
      Body: body,
   }
   post.Timestamp = post.PostID.Time().Format(timestampFormat)
   post.TimestampShort = post.PostID.Time().Format(timestampShortFormat)
   err := post.WritePost()
   if err != nil {
      return nil, err
   }
   return post, nil
}

func (post *Post) WritePost() error {
   contents, err := json.Marshal(post)
   if err != nil {
      return err
   }
   path := filepath.Join(docDir, post.PostID.String() + ".json")
   return ioutil.WriteFile(path, contents, 0644)
}

func LoadPost(id string) (*Post, error) {
   post := &Post{}
   path := filepath.Join(docDir, id + ".json")
   contents, err := ioutil.ReadFile(path)
   if err != nil {
      return nil, err
   }
   err = json.Unmarshal(contents, post)
   if err != nil{
      return nil, err
   }
   post.PostID, err = ksuid.Parse(id)
   if err != nil{
      return nil, err
   }
   post.Timestamp = post.PostID.Time().Format(timestampFormat)
   post.TimestampShort = post.PostID.Time().Format(timestampShortFormat)
   return post, nil
}

type Posts struct {
   PostSlice []string
   PostMap   map[string]*Post
}

func (posts *Posts) AddPost(post *Post) {
   id := post.PostID.String()
   posts.PostMap[id] = post
   posts.PostSlice = append(posts.PostSlice, id)
}

func LoadPosts() (*Posts, error) {
   posts := &Posts{}
   posts.PostMap = make(map[string]*Post)
   files, err := ioutil.ReadDir(docDir)
   if os.IsNotExist(err) {
      err := os.Mkdir(docDir, 0644)
      if err != nil {
         return nil, err
      }
      return posts, nil
   } else if err != nil {
      return nil, err
   }
   for _, file := range files {
      if file.IsDir() { continue }
      filename := file.Name()
      extension := filepath.Ext(filename)
      if extension == ".json" {
         id := strings.TrimSuffix(filename, extension)
         post, err := LoadPost(id)
         if err != nil {
            return nil, err
         }
         posts.AddPost(post)
      }
   }
   return posts, nil
}

// Pretty much a more barebones Post object.
type Announcement struct {
   AnnouncementID ksuid.KSUID
   Timestamp      string
   TimestampShort string
   Body           string
}

func PublishAnnouncement(body string) (*Announcement, error) {
   announcement := &Announcement{}
   announcement.Body = body
   announcement.AnnouncementID = ksuid.New()
   announcement.Timestamp = announcement.AnnouncementID.Time().Format(timestampFormat)
   announcement.TimestampShort = announcement.AnnouncementID.Time().Format(timestampShortFormat)
   path := filepath.Join(announcementDir, announcement.AnnouncementID.String() + ".json")
   err := ioutil.WriteFile(path, []byte(body), 0644)
   if err != nil {
      return nil, err
   }
   return announcement, nil
}

func LoadAnnouncement(id string) (*Announcement, error) {
   announcement := &Announcement{}
   path := filepath.Join(announcementDir, id + ".json")
   contents, err := ioutil.ReadFile(path)
   if err != nil {
      return nil, err
   }
   announcement.Body = string(contents)
   announcement.AnnouncementID, err = ksuid.Parse(id)
   if err != nil {
      return nil, err
   }
   announcement.Timestamp = announcement.AnnouncementID.Time().Format(timestampFormat)
   announcement.TimestampShort = announcement.AnnouncementID.Time().Format(timestampShortFormat)
   return announcement, nil
}

type Announcements struct {
   AnnouncementSlice []string
   AnnouncementMap   map[string]*Announcement
}

func (announcements *Announcements) AddAnnouncement(announcement *Announcement) {
   id := announcement.AnnouncementID.String()
   announcements.AnnouncementMap[id] = announcement
   announcements.AnnouncementSlice = append(announcements.AnnouncementSlice, id)
}

func LoadAnnouncements() (*Announcements, error) {
   announcements := &Announcements{}
   announcements.AnnouncementMap = make(map[string]*Announcement)
   files, err := ioutil.ReadDir(announcementDir)
   if os.IsNotExist(err) {
      err := os.Mkdir(announcementDir, 0644)
      if err != nil {
         return nil, err
      }
      return announcements, nil
   } else if err != nil {
      return nil, err
   }
   for _, file := range files {
      if file.IsDir() { continue }
      filename := file.Name()
      extension := filepath.Ext(filename)
      if extension == ".json" {
         id := strings.TrimSuffix(filename, extension)
         announcement, err := LoadAnnouncement(id)
         if err != nil {
            return nil, err
         }
         announcements.AddAnnouncement(announcement)
      }
   }
   return announcements, nil
}
