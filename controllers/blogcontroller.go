package controllers

import "log"

type BlogController struct {
}

func (this *BlogController) init() bool {
  tools := registerMenuCategory('Tools')
  
  registerMenuOption(&InitiumMenu{title: 'Home', 
    route: registerRoute(&InitiumRoute{uri: '/', callback: this.Index}),
    category: tools,
  })

  registerMenuOption(&InitiumMenu{title: 'News',
    route: registerRoute(&InitiumRoute{uri: '/news', callback: this.News),
  })

  registerRoute(&InitiumRoute{uri: '/news', callback: this.CreateNews, method: MethodPOST})

  return true
}

func (this *BlogController) Index(a bool, b bool) error {
  log.Println("Index from blog controller.")
  return nil
}

func (this *BlogController) News(a bool, b bool) error {
  return nil
}

func (this *BlogController) CreateNews(a bool, b bool) error {
  return nil
}