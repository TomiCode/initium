/* 
      Initium Project 
  ** Tomasz Król, 2016 **
*/

;(function($, window, document, undefined) {

"use strict";

var window = (typeof window != 'undefined' && window.Math == Math)
  ? window : (typeof self != 'undefined' && self.Math == Math)
    ? self : Function('return this')();

$.iforms = $.fn.iforms = function(parameters) {

  $(this).each(function(){
    var
      settings = (parameters !== undefined && $.isPlainObject($parameters)) 
        ? $.extend(true, {}, $.fn.iforms.settings, parameters) : $.fn.iforms.settings,
      element = this,
      $module = $(this),
      module;

    module = {
      initialize: function() {
        console.log("Initialized initium.forms:", $module);
        module.attach();
      },
      destroy: function() {
        console.log("Destroyed for", $module);
      },
      attach: function() {
        console.log("Attaching events for module: ", $module);
        var
          event = module.get.event();

        if(event) {
          $module.on(event, module.event);
        }
      },
      event: function(event) {
        module.handle();
        if(event.type == 'submit' || event.type == 'click') {
          event.preventDefault();
        }
      },
      handle: function() { 
        console.log("Handle module request.");
        module.send();
        module.remove.message();
      },
      state: {
        loading: function() {
          return (module.xhr) ? (module.xhr.state() == 'pending') : false;
        },
        disabled: function() {
          return $module.hasClass("disabled");
        }
      },
      get: {
        event: function() {
          if($module.is('input')) {
            console.log("Invalid object to handle events:", $module);
            return false;
          }
          else if($module.is('form')) {
            return 'submit';
          }
          else {
            return 'click';
          }
        },
        routing: function(data) {
          var
            handler = $module.data('handler')  || false,
            address = settings.routes[handler] || false,
            handlerParameters;

          if(address) {
            handlerParameters = address.match(/\{\$*[A-Za-z0-9]+\}/g);
            if(handlerParameters) {
              console.log("Handling required routing parameters:", handlerParameters);
              $.each(handlerParameters, function(index, template){
                var
                  variable = template.substr(1, template.length - 2),
                  value = ($.isPlainObject(data) && data[variable] !== undefined) ? data[variable] :
                    ($module.data(variable) !== undefined) ? $module.data(variable) : undefined;

                if(value === undefined) {
                  console.log("Can not parse required parameter:", variable);
                  return false;
                }
                else {
                  console.log("Found parameter:", variable, "with value:", value);
                  address = address.replace(template, value);
                }
              });
            }
          }
          return address;
        },
        content: function() {
          if($module.is('form')) {
            return $module.serialize();
          }
          return undefined;
        }
      },
      send: function() {
        if(!module.state.loading()) {
          console.log("Creating xhr reqiest.");
          module.set.loading();
          module.xhr = module.request.create();
        } 
        else {
          console.log("Existent xhr request is already pending!");
        }
      },
      request: {
        create: function() {
          var
            data = module.get.content(),
            url  = module.get.routing(),
            xhr;

          console.log("Creating xhr request to:", url);
          if(url) {
            xhr = $.ajax({
              url       : url,
              data      : data,
              method    : "POST",
              // completed : module.request.completed,
              // error     : module.request.error,
              // success   : module.request.success,
            })
              .always(module.request.always)
              .done(module.request.done)
              .fail(module.request.fail);
          }
          return xhr;
        },
        always: function(res, status, obj) {
          console.log("Ajax completed:", status);
          setTimeout(function(){
            module.remove.loading();
          }, 500);
        },
        done: function(data, status, xhr) {
          console.log("Ajax success:", status, "data:", data);
          if(data.success !== undefined && !data.success) {
            if(data.error !== undefined) {
              module.set.message.error(data.error);
            }
            else {
              module.set.message.error("An error occurred. Please try again.");
            }
          }
        },
        fail: function(xhr, status, error) {
          console.log("Ajax failed:", status, "error:", error);
        }
      },
      set: {
        loading: function() {
          console.log("Adding loading class into object..");
          $module.addClass("loading");
        },
        disabled: function() {
          console.log("Adding disabled class into object..");
          $module.addClass("disabled");
        },
        error: function() {
          console.log("Adding error class into object..");
          $module.addClass("error");
        },
        message: {
          error: function(content) {
            console.log("Adding error message into module:", content);
            if($module.is('form')) {
              console.log("Module is a form. Creating error message element.");
              $module.append(
                $('<div class="ui message error">')
                  .html(content)
                  .transition("fade"));
            }
          }
        }
      },
      remove: {
        loading: function() {
          console.log("Removing loading class..");
          $module.removeClass("loading");
        },
        disabled: function() {
          console.log("Removing disabled class..");
          $module.removeClass("disabled");
        },
        error: function() {
          console.log("Removing error class..");
          $module.removeClass("error");
        },
        message: function() {
          if($module.is('form')) {
            var $message = $module.find(".ui.error");
            if ($message !== undefined) {
              $message.transition('fade');
            }
          }
        }
      }
    };

    module.initialize();
  });
  return this;
};

$.iforms.settings = {

  routes: {},

};

})( jQuery, window, document );
