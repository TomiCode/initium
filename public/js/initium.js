$.fn.iforms.settings.routes = {
  'auth.login.form': "/auth/form",
  'test.abc': '/test/{userid}',
};

$.fn.api.settings.api = {
  "auth.login.form": "/auth/form",
};

$(document).ready(function(){
  $('.secondary.menu').visibility({
    type: 'fixed'
  });

  $('.ui.dropdown').dropdown();
});