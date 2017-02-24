/* 
      Initium Project 
  ** Tomasz Król, 2016 **
*/

;(function($, window, document, undefined) {

"use strict";

var window = (typeof window != 'undefined' && window.Math == Math)
  ? window : (typeof self != 'undefined' && self.Math == Math)
    ? self : Function('return this')();


$.imsg = $.fn.imsg = function(parameters) {
  $(this).each(function() {
    var
      $module = $(this),
      element = this,
      module;

    module = {
      init: function() {
        if($module.css('display') == 'none') {
          $module.transition('fade');
        }
        $module.on('mouseenter', module.mouse.enter);
        $module.on('mouseleave', module.mouse.leave);

        module.timer = setTimeout(module.destroy, 3500);
      },
      started: function() {
        return module.timer || false;
      },
      animate: function() {
        return module.animating || false;
      },
      mouse: {
        enter: function() {
          clearTimeout(module.timer);
          module.timer = false;
          $module.transition('pulse');
        },
        leave: function() {
          if(!module.timer) {
            module.timer = setTimeout(module.destroy, 2000);
          }
        }
      },
      destroy: function() {
        $module.transition('fade', function(){
          $module.remove();
        });
      }
    }
    module.init();
  });
}

$.iforms = $.fn.iforms = function(parameters) {

  $(this).each(function(){
    var
      settings = (parameters !== undefined && $.isPlainObject($parameters)) 
        ? $.extend(true, {}, $.fn.iforms.settings, parameters) : $.fn.iforms.settings,

      element = this,
      $module = $(this),
      module,

      requestTime,
      bindType = $module.data("binding") || false,
      $context = (bindType && bindType == 'form') ? $module.closest('form') : $module,

      eventSuffix = '.' + settings.namespace + '.iforms'
    ;

    if($context.length == 0) {
      console.warn("Invalid context for element:", $module);
      $context = $module;
    }

    module = {
      initialize: function() {
        console.log("Initialized initium.forms:", $module, "binding type:", bindType);
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
          $module.on(event + eventSuffix, module.event);
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
        if(module.validate()) {
          module.set.loading();
          module.send();
        }
        else {
          console.log("Validator errors, fix inserted values!");
        }
      },
      state: {
        loading: function() {
          return (module.xhr) ? (module.xhr.state() == 'pending') : false;
        },
        disabled: function() {
          return $module.hasClass("disabled");
        }
      },
      is: {
        form: function() {
          return (bindType && bindType == 'form');
        },
      },
      validate: function() {
        if(!module.is.form()) {
          return true;
        }

        var
          $inputs = $context.find('input[data-validator]'),
          $fields = $context.find('.field'),
          fieldsValid = true
        ;

        $inputs.each(function(){
          var
            $field = $(this),
            $fieldGroup = $field.closest($fields),
            type = $field.data('validator'),
            validator = $.iforms.validators[type],
            field,

            $prompt
          ;

          field = {
            validate: function() {
              if(!field.test()) {
                $field.on('input.initium', field.change);
                $fieldGroup.addClass("error");
                if(validator.error !== undefined) {
                  field.create.prompt(validator.error);
                }
                else {
                  field.create.prompt("Invalid value!");
                }

                field.failed = true;
                fieldsValid = false;
              }
            },
            test: function() {
              var 
                count = $field.val().length;

              if(validator !== undefined) {
                if(validator.min !== undefined && validator.min > count) {
                  return false;
                }
                if(validator.max !== undefined && validator.max < count) {
                  return false;
                }
                if(validator.expr !== undefined) {
                  return validator.expr.test($field.val());
                }
              }
              console.log("Unknown validator:", type);
              return false;
            },
            error: function() {
              return field.failed || false;
            },
            change: function() {
              console.log("Error input onChange called.");
              $field.off('input.initium');
              if(field.error()) {
                $fieldGroup.removeClass('error');
                field.failed = false;

                $prompt.transition("scale out", function(){
                  $prompt.remove();
                });
              }
            },
            template: {
              prompt: function(message) {
                return $('<div/>').addClass('ui basic red pointing prompt label').html(message);
              }
            },
            create: {
              prompt: function(msg) {
                $prompt = field.template.prompt(msg);
                $prompt.appendTo($fieldGroup);

                $prompt.transition('scale in');
              }
            }
          }
          field.validate();
        });
        return fieldsValid;
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
            handler = $context.data('handler') || false,
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
          if(module.is.form()) {
            return $context.serialize();
          }
          return undefined;
        }
      },
      send: function() {
        if(!module.state.loading()) {
          console.log("Creating xhr reqiest.");
          module.xhr = module.request.create();
        } 
        else {
          console.log("Existent xhr request is already pending!");
        }
      },
      response: {
        handle: function(data) {
          module.remove.loading();
          if(data.success === false) {
            $module.addClass('disabled');
            if(data.error !== undefined) {
              module.set.message.error(data.error);
            }
            else {
              module.set.message.error("An error occurred. Please try again.");
            }
            setTimeout(function(){
              $module.removeClass('disabled');
              // module.remove.message();
            }, (data.error !== undefined) ? 2100 : 1201 );
          }
          else if(data.success === true) {
            module.remove.message();
            $module.addClass("positive");
            
            if(data.redirect !== undefined) {
              $module.addClass("loading");
              setTimeout(function(){
                window.location.href = data.redirect;
              }, 500);
            }
          }
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
            })
              .always(module.request.always)
              .done(module.request.done)
              .fail(module.request.fail);
          }
          return xhr;
        },
        always: function(res, status, obj) {
          console.log("Ajax completed:", status);
          // module.remove.loading();
          // setTimeout(function(){
            
          // }, 500);
        },
        done: function(data, status, xhr) {
          var 
            delay = (settings.delay - (new Date().getTime() - requestTime));

          delay = (delay > 0 ? delay : 0);

          console.log("Ajax success:", status, "data:", data, "delay:", delay);
          setTimeout(function(){ module.response.handle(data); }, delay);
        },
        fail: function(xhr, status, error) {
          console.log("Ajax failed:", status, "error:", error);
        }
      },
      set: {
        loading: function() {
          console.log("Adding loading class into object..");
          requestTime = new Date().getTime();
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
            module.message = true;

            if(module.is.form()) {
              console.log("Module is a form. Creating error message element.");
              
              var $message = $context.find('.ui.error');
              if ($message.length == 0) {
                $message = $('<div class="ui message error" style="display: none;">');
                $context.prepend($message);
                
              }
              $message.html(content).transition('fade in');
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
          if(module.message) {
            if(module.is.form()) {
              var $message = $context.find(".ui.error");
              if ($message !== undefined) {
                $message.transition('fade');
              }
            }
            module.message = false;
            return true;
          }
          return false;
        }
      }
    };

    module.initialize();
  });
  return this;
};

$.iforms.settings = {
  namespace: "initium",
  routes: {},
  delay: 1000,
};

$.iforms.validators = {
  email: {
    min: 3,
    expr: /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w+)+$/,
    error: "Invalid email address"
  }
}

})( jQuery, window, document );
