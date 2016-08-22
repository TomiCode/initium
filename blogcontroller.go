package main

import "fmt"

type BlogController struct {
  App ApplicationInterface
}

type BlogPost struct {
  Title string
  Content string

  View int
  Like int
}

func (controller* BlogController) RoutingRegister() []ControllerRoute {
  return []ControllerRoute{
    ControllerRoute{uri: "/", call: controller.index},
    ControllerRoute{uri: "/add/{user}", call: controller.addPost},
  }
}

func (controller* BlogController) index(req *InitiumRequest) error {
  fmt.Println("Index sees params:", req.vars)

  var test_posts = []BlogPost{
    BlogPost{Title: "First blog entry 01", Content: "Lorem ipsum.", View: 1337, Like: 0},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
  }

  return controller.App.RenderTemplate(req, "blog.index", test_posts)
}

func (controller* BlogController) addPost(req *InitiumRequest) error {
  fmt.Println("BlogController addPost.")
  return nil
}
