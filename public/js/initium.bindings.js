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

/* jQuery function to bind buttons with the Initium Handler Module. */
$.initium.bind = $.fn.initium.bind = function(parameters) {
  $(this).each(function() {
    var
      core = $.initium.core,
      settings = (parameters !== undefined && $.isPlainObject(parameters)) ?
        $.extend(true, {}, $.initium.settings, parameters) : $.initium.settings,

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
$.initium.routing = {};

/* Binding default settings. */
$.initium.settings = {
  
  debug: true,

};

})( jQuery, window, document );