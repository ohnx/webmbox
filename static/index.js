/* global Vue, app */

const API_LOCATION = 'http://smtp.masonx.ca/';

Vue.component('email-small', {
  props: ['todo'],
  template: '<li>{{ todo.text }}</li>'
});

var listMailboxes = function() {
  var xmlHttp = new XMLHttpRequest();

  xmlHttp.onreadystatechange = function() { 
      if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
        app.mailboxes = JSON.parse(xmlHttp.responseText);
      }
  };

  xmlHttp.open('GET', API_LOCATION, true); // true for asynchronous 
  xmlHttp.send(null);
};

var openMailbox = function(mailbox) {
  var xmlHttp = new XMLHttpRequest();

  xmlHttp.onreadystatechange = function() { 
      if (xmlHttp.readyState == 4) {
        if (xmlHttp.status == 200) {
          app.emails = JSON.parse(xmlHttp.responseText);
        } else {
          document.querySelector('#toast').MaterialSnackbar.showSnackbar(JSON.parse(xmlHttp.responseText));
        }
      }
  };

  xmlHttp.open('GET', API_LOCATION+mailbox, true); // true for asynchronous 
  xmlHttp.send(null);
};

(function() {
  window.app = new Vue({
    el: '#app',
    data: {
      emails: [],
      mailboxes: [],
      mailbox: ''
    },
    methods: {
      openmb: function(mailbox) {
        listMailboxes();
        this.mailbox = mailbox;
        openMailbox(this.mailbox);
      },
      expand: function(index) {
        if (this.emails[index].Content.length > 1) {
          /* has html-formatted content */
          this.emails[index].Content[0].Data = this.emails[index].Content[1].Data;
        }
      }
    }
  });

  console.log(window.app);
  listMailboxes();
})();
