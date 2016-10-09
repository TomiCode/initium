/***
 *
 *  Developed and Programmed by Tomasz Król <tomicode@gmail.com>, 
 *   INITIUM PROJECT (c) 2016.
 *    Afterall, if You read this message, then have a nice day! :)
 *
 * Javascript bindings for Initium content javascript bindings 
 *  and those nasty animations.
 *
 ***/

 ;(function($, window, document, undefined) {

"use strict";

var window = (typeof window != 'undefined' && window.Math == Math)
  ? window : (typeof self != 'undefined' && self.Math == Math)
    ? self : Function('return this')();

/* Project namespace. */
$.initium = $.fn.initium = {
  core : {
    log: function() {
      if($.initium.settings.debug) {
        console.log.apply(console, arguments);
      }
    },
    warn: function() {
      console.warn.apply(console, arguments);
    },
    error: function() {
      console.error.apply(console, arguments);
    }
  }
};

$.initium.validate = $.fn.initium.validate = function(parameters) {
  var
    core = $.initium.core,
    settings = (parameters !== undefined && $.isPlainObject(parameters)) ?
        $.extend(true, {}, $.initium.settings, parameters) : $.initium.settings,

    $fields = $(this),
    success = true;

  $fields.each(function(){
    var
      element = this,
      $field = $(this),
      field,

      validator = $.initium.validators[$field.data(settings.data.validator)] || false;

    field = {
      initialize: function() {
        if(!validator) {
          core.warn("Validator not declared for this field:", $field);
          return;
        }
        core.log("Initialized field verification:", $field);
      },
      input: function() {

      },
      checkbox: function() {

      },
      dropdown: function() {

      },
      verificate: function() {
        if(field.input()) {

        }
        else if(field.checkbox()) {

        }
        else if(field.dropdown()) {

        }
      },
      error: function() {
        core.log("Error occurred while field validation.");
        if(success) {
          success = false;
        }
      }
    };
    field.initialize();
  });

  /* Return the validator result. */
  return success;
};

/* jQuery function to bind buttons with the Initium Handler Module. */
$.initium.bind = $.fn.initium.bind = function(parameters) {
  var 
    core = $.initium.core,
    settings = (parameters !== undefined && $.isPlainObject(parameters)) ?
        $.extend(true, {}, $.initium.settings, parameters) : $.initium.settings,

    $modules = $(this);

  $modules.each(function() {
    var
      element = this,
      $module = $(this),
      module;

    module = {
      initialize: function() {
        core.log("Hello! Initialized Initium binding for:", $module);
      },
      destroy: function() {
        core.warn("Destroyed Initium binding for:", $module);
      },
    };
    module.initialize();
  });
  return this;
};

/* Initium Handler routing table. */
$.initium.routes = {};

/* Field validators. */
$.initium.validators = {};

/* Binding default settings. */
$.initium.settings = {
  debug: true,

  /* Validator type field. */
  data: {
    validator: "validator"
  }
};

})( jQuery, window, document );