package controllers

import "log"

type BlogController struct {
  BasicController
}

func (this *BlogController) alias() string {
  return "blog"
}

func (this *BlogController) init() bool {
  log.Println("Registering BlogController into Initium..")

  tools := registerCategory(this.id, "Tools")

  // 'blog'
  registerOption(this.id, "Home", BaseCategory,
    registerRoute(this.id, "", this.Index, MethodGET, nil))

  // 'blog_news'
  registerOption(this.id, "News", tools,
    registerRoute(this.id, "/news", this.News, MethodGET, nil))

  // 'blog_news'  
  registerRoute(this.id, "/news", this.CreateNews, MethodPOST, nil)

  // 'blog_news_id'
  registerRoute(this.id, "/news/:id", this.ViewNews, MethodGET, nil)
  // 'blog_news_id'
  registerRoute(this.id, "/news/:id", this.EditNews, MethodPOST, nil)

  // 'blog_news_id_delete'
  registerRoute(this.id, "/news/:id/delete", this.DeleteNews, MethodGET, nil)

  return true
}

func (this *BlogController) Index(a bool, b bool) error {
  log.Println("BlogController Index.")
  return nil
}

func (this *BlogController) News(a bool, b bool) error {
  log.Println("BlogController News.")
  return nil
}

func (this *BlogController) CreateNews(a bool, b bool) error {
  log.Println("BlogController CreateNews.")
  return nil
}

func (this *BlogController) ViewNews(a bool, b bool) error {
  log.Println("BlogController ViewNews.")
  return nil
}

func (this *BlogController) EditNews(a bool, b bool) error {
  log.Println("BlogController EditNews.")
  return nil
}

func (this *BlogController) DeleteNews(a bool, b bool) error {
  log.Println("BlogController DeleteNews.")
  return nil
}
