/* 
      Initium Project 
  ** Tomasz Król, 2016 **
*/

;(function($, window, document, undefined) {

"use strict";

var window = (typeof window != 'undefined' && window.Math == Math)
  ? window : (typeof self != 'undefined' && self.Math == Math)
    ? self : Function('return this')();

$.fapi = $.fn.fapi = function(parameters) {
  // var $modules = $.isFunction(this) ? $(window) : $(this);

  $(this).each(function(){
    var
      element = this,
      $module = $(this),
      module;

    module = {
      initialize: function() {
        console.log("Initialized with jQuery module", $module);
        module.bind();
      },
      destroy: function() {
        console.log("Destroyed for", $module);
      },
      bind: function() {
        console.log("Binding event to control.");
        var event = module.local.event();
        if (event) {
          $module.on(event, module.self.event);
        }
      },
      local: {
        event: function() {
          if ($.isWindow(element)) {
            console.log("fapi called on Window. Can not attach to events.");
            return false;
          }
          if ($module.is('input')) {
            console.log("fapi is not for input elements!");
            return false;
          }
          else if ($module.is('form')) {
            return 'submit';
          }
          else {
            return 'click';
          }
        },
        route: function(handler) {
          if (initium.routes !== undefined && handler !== undefined) {
            var 
              route = initium.routes[handler]

            if (route !== undefined) {
              return route;
            }
          }
          return false;
        },
        data: function() {
          if ($module.is('form')) {
            return $module.serialize();
          }
          else {
            return {nodata: true};
          }
        },
        request: function(url) {
          console.log("Sending request to:", url);

          var 
            data = module.local.data();

          console.log("Request data:", data);
          $.ajax({ 
            method: "POST",
            url: url,
            data: data
          })
            .done(module.self.request.done)
            .fail(module.self.request.fail)
            .always(module.self.request.always);
        }
      },
      self: {
        event: function(ev) {
          ev.preventDefault();
          if (module.self.state.loading()) {
            return;
          }

          module.set.loading(true);

          var
            handler = $module.data("handler"),
            route = module.local.route(handler);

          if (route) {
            module.local.request(route);
          } 
          else {
            module.set.loading(false);
          }
        },
        state: {
          loading: function() {
            return module.loading || false;
          },
          disabled: function() {
            return module.disabled || false;
          }
        },
        request: {
          done: function(data, status, xhr) {
            console.log("XHRRequest done. ", data);
          },
          fail: function(xhr, status, error) {
            console.log("XHR failed: ", error);
          },
          always: function(res, status, obj) {
            console.log("XHR always: ", status);
            setTimeout(function(){
              module.set.loading(false);
            }, 500);
          }
        }
      },
      set: {
        loading: function(b) {
          if (!module.self.state.loading() && b) {
            $module.addClass("loading");
            module.loading = true;
          }
          else if (module.self.state.loading() && !b) {
            $module.removeClass("loading");
            module.loading = false;
          }
        },
        disabled: function(b) {
          if (!module.self.state.disabled() && b) {
            $module.addClass("disabled");
            module.disabled = true;
          }
          else if(module.self.state.disabled() && !b) {
            $module.removeClass("disabled");
            module.disabled = false;
          }
        }
      }
    };

    module.initialize();
  });
  
  return this;
}

})( jQuery, window, document );
