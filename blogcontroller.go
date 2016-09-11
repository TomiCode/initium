package main

import "log"

type BlogController struct {
  App ApplicationInterface
}

type BlogPost struct {
  Title string
  Content string
  View int
  Like int
}

func (controller *BlogController) RegisterModule() *InitiumModule {
  return &InitiumModule{Title: "Blog", RouteName: "blog.index", ControllerAlias: "blog_ctrl"}
}

func (controller *BlogController) RegisterOptions() []*InitiumModuleCategory {
  return []*InitiumModuleCategory{
    &InitiumModuleCategory{
      Title: "", 
      Options: []*ModuleOption{
        &ModuleOption{Name: "Add entry", Route: "blog.add"},
      },
    },
  }
}

func (controller *BlogController) RegisterRouting() []*ControllerRoute {
  return []*ControllerRoute{
    &ControllerRoute{uri: "/", call: controller.index, alias: "blog.index"},
    &ControllerRoute{uri: "/add", call: controller.addPost, alias: "blog.add", access: Permission_Auth_None},
  }
}

func (controller *BlogController) index(req *InitiumRequest, params *RequestParameters) error {
  log.Println("[BlogController] Parameters:", params)
  
  var test_posts = []BlogPost{
    BlogPost{Title: "First blog entry 01", Content: "Lorem ipsum.", View: 1337, Like: 0},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
    BlogPost{Title: "Testing golang templating systems", Content: "Lorem ipsum lorem ipsum lorem ipsum", View: 3},
  }
  
  return controller.App.RenderTemplate(req, "blog.index", test_posts)
}

func (controller *BlogController) addPost(req *InitiumRequest, params *RequestParameters) error {
  return controller.App.RenderTemplate(req, "blog.add", nil)
}
