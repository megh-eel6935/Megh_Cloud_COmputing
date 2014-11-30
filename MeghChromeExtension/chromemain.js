  var currentgroupidupload = "usertime";
  var currentgroupidview = "usertime";
  var currentgroupname = 0;
  var urlRegex = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
  var ChatRoom = function() {
      var c = this;

      c.create = function() {
         
          return c;
      };

      c.connect = function() {

          console.log("ws://" + location.host + "/sockets/sanath");
          c.ws = new WebSocket("ws://" + location.host + "/sockets/sanath");
          return c;
      };

      c.disconnect = function() {
          c.ws.close();
      };

      c.listen = function() {
          c.ws.addEventListener("message", function(e) {
              console.log("received" + e.data);

              var temp = JSON.parse(e.data);

              if (currentgroupid == temp.groupid) {
                  updatelist(temp.groupid);
              } else {
                  val = $("#badge" + temp.groupid).text();
                  if (val == undefined) {
                      val = 0;
                  }
                  val = val + 1;
                  $("#badge" + temp.groupid).text(val);
              }


          });

          c.ws.addEventListener("close", function(e) {
              console.log("closed connection, reconnecting", e)
              c.connect()
          });

          c.ws.addEventListener("error", function(e) {
              console.log("error for connection", e)
          });

          return c;
      };

      c.send = function(messageData) {
          console.log("sent" + JSON.stringify(messageData));
          c.ws.send(JSON.stringify(messageData));
      };

      c.handleSend = function() {
          
      };

      var initialize = function() {
          c.create().handleSend();
          c.connect().listen();
      }();
  };

  $(document).ready(function() {
      chat = new ChatRoom();

      $.get("/groupslist", function(data) {
          sap = JSON.stringify(data);
          console.log(data);
          var length = data.data.length;
          for (var i = 0; i < length; i++) {

              if (data.data[i].username == data.data[i].group_admin) {
                  var temp = "<option id=" + data.data[i].group_id + " value=" + data.data[i].group_id + ">" + data.data[i].group_name + "</option>";

              } else {
                  var temp = "<option id=" + data.data[i].group_id + " value=" + data.data[i].group_id + ">" + data.data[i].group_name + "</option>";
              }

              $("#groupslistupload").append(temp);
              $("#groupslistview").append(temp);


          }

      });


      $.get("/getmessages", function(data) {
          sap = JSON.stringify(data);
          var length = data.data.length;
          console.log("from userdata" + length);
          for (var i = 0; i < length; i++) {
            if (data.data[i].content_type == "text") {
                      var timeline = "<li id=\"li" + data.data[i].timestamp + "\"><div class=\"timeline-badge warning\"><i class=\"glyphicon glyphicon-envelope\"></i></div><div class=\"timeline-panel\">  <div id=" + data.data[i].timestamp + " class=\"close\">&times;</div>" + "<div class=\"timeline-heading\"><h4 class=\"timeline-title\">" + linkify(data.data[i].content) + "</h4><p><small class=\"text-muted\">" + "<i class=\"glyphicon glyphicon-time\"></i>" + timeConverter(data.data[i].timestamp) + "</small></p></div></div></li>"

                  } else if (data.data[i].content_type == "url") {
                      var timeline = "<li id=\"li" + data.data[i].timestamp + "\"><div class=\"timeline-badge info\"><i class=\"glyphicon glyphicon-link\"></i></div><div class=\"timeline-panel\"> <div id=" + data.data[i].timestamp + " class=\"close\">&times;</div>" + "<div class=\"timeline-heading\"><h4 class=\"timeline-title\">" + linkify(data.data[i].content) + "</h4><p><small class=\"text-muted\">" + "<i class=\"glyphicon glyphicon-time\"></i>" + timeConverter(data.data[i].timestamp) + "</small></p></div></div></li>"
                  } else if (data.data[i].content_type == "file") {

                      var timeline = "<li id=\"li" + data.data[i].timestamp + "\"><div class=\"timeline-badge danger\"><i class=\"glyphicon glyphicon-file\"></i></div><div class=\"timeline-panel\">   <div id=" + data.data[i].timestamp + " class=\"close\">&times;</div>" + "<div class=\"timeline-heading\"><h4 class=\"timeline-title\"><a href='https://s3-us-west-2.amazonaws.com/megh-uploads/" + data.data[i].content + "' class=\"glyphicon glyphicon-download-alt\">  </a> " + data.data[i].content + "</h4><p><small class=\"text-muted\">" + "<i class=\"glyphicon glyphicon-time\"></i>" + timeConverter(data.data[i].timestamp) + "</small></p></div></div></li>"

                  }
                  $("#usertimeline").append(timeline);
          
          }



      });

      $("#groupslistupload").click(function() {

          currentgroupidupload = $("#groupslistupload").val();

          console.log("currentgroupid is " + currentgroupidupload);

      });

      $("#groupslistview").click(function() {

          currentgroupidview = $("#groupslistview").val();

          id = currentgroupidview;
          console.log("currentgroupid is " + currentgroupidview);
          if (id == "usertime") {
              var url = "/getmessages";

          } else {
              var url = "/getgroupdatabyid/" + id;
              var url2 = "/usersingroup/" + id;
          }
          $.get(url, function(data) {
              sap = JSON.stringify(data);
              var length = data.data.length;
              console.log("from getgroupdata " + length);
              $("#usertimeline").empty();
              $("#userslist").empty();
              for (var i = 0; i < length; i++) {
                  if (data.data[i].content_type == "text") {
                      var timeline = "<li id=\"li" + data.data[i].timestamp + "\"><div class=\"timeline-badge warning\"><i class=\"glyphicon glyphicon-envelope\"></i></div><div class=\"timeline-panel\">  <div id=" + data.data[i].timestamp + " class=\"close\">&times;</div>" + "<div class=\"timeline-heading\"><h4 class=\"timeline-title\">" + linkify(data.data[i].content) + "</h4><p><small class=\"text-muted\">" + "<i class=\"glyphicon glyphicon-time\"></i>" + timeConverter(data.data[i].timestamp) + "</small></p></div></div></li>"

                  } else if (data.data[i].content_type == "url") {
                      var timeline = "<li id=\"li" + data.data[i].timestamp + "\"><div class=\"timeline-badge info\"><i class=\"glyphicon glyphicon-link\"></i></div><div class=\"timeline-panel\"> <div id=" + data.data[i].timestamp + " class=\"close\">&times;</div>" + "<div class=\"timeline-heading\"><h4 class=\"timeline-title\">" + linkify(data.data[i].content) + "</h4><p><small class=\"text-muted\">" + "<i class=\"glyphicon glyphicon-time\"></i>" + timeConverter(data.data[i].timestamp) + "</small></p></div></div></li>"
                  } else if (data.data[i].content_type == "file") {

                      var timeline = "<li id=\"li" + data.data[i].timestamp + "\"><div class=\"timeline-badge danger\"><i class=\"glyphicon glyphicon-file\"></i></div><div class=\"timeline-panel\">   <div id=" + data.data[i].timestamp + " class=\"close\">&times;</div>" + "<div class=\"timeline-heading\"><h4 class=\"timeline-title\"><a href='https://s3-us-west-2.amazonaws.com/megh-uploads/" + data.data[i].content + "' class=\"glyphicon glyphicon-download-alt\">  </a> " + data.data[i].content + "</h4><p><small class=\"text-muted\">" + "<i class=\"glyphicon glyphicon-time\"></i>" + timeConverter(data.data[i].timestamp) + "</small></p></div></div></li>"

                  }
                  $("#usertimeline").append(timeline);
              }

          });

      });


      $(document).on('click', '.close', function() {

          if (confirm('Are  You Sure you Want to delete tis item')) {
              var id = this.id;
              if (currentgroupidview == "usertime") {
                  var url = '/deleteuserdata/' + id;
                  $.get(url, function(data) {
                      console.log(data.status);
                      if (data.status == "Success") {
                          console.log(data.status);
                          $('#li' + id).fadeOut(500, function() {
                              $('#li' + id).remove();
                          });
                      }
                  });
              } else {
                  var url = '/deletegroupdata/' + currentgroupidview + '/' + id;
                  $.get(url, function(data) {
                      console.log(data.status);
                      if (data.status == "Success") {

                          $('#li' + id).fadeOut(500, function() {
                              $('#li' + id).remove();
                          });
                      }
                  });
              }
          } else {}

      });


      $("#userfilesend").click(function() {
          var value = $("#text-input").val();
          var id = currentgroupidupload;

          if (checkurl(value)) {
              var isurl = true;
              if (id == "usertime") {
                  var url = "/urls";
              } else {
                  var url = "/groupdataurl/" + id;
              }
          } else {
              if (id == "usertime") {
                  var url = "/messages";
              } else {
                  var url = "/groupdata/" + id;
              }

          }
          $.post(url, {
              message: value
          }, function(data, status) {
             if(data.status=="success"){
             $("#text-input").val("");
             $.jGrowl("Successfull");
             $('#groupslistview').val(''+currentgroupidupload);
              updatelist(id);
              if (isurl) {
                  chat.ws.send('{"type":"url","from":"' + getCookie(name) + '","groupid":"' + currentgroupidupload + '","content":"' + value + '"}');

              } else {

                  chat.ws.send('{"type":"text","from":"' + getCookie(name) + '","groupid":"' + currentgroupidupload + '","content":"' + value + '"}');
              }
            }
          });
      });

      $("#userfile").change(function() {
          var id = currentgroupidupload;

          if (id == "usertime") {
              var url = "/uploadfile";
          } else {
              var url = "/groupuploadfile/" + id;
          }
          var filename = $("#userfile").val().split('\\').pop();
          console.log("filename is " + filename);

          var file = $("#userfile")[0].files[0];
          var formdata = new FormData();
          formdata.append("file", file);
          $.ajax({
              url: url,
              data: formdata,
              processData: false,
              contentType: false,
              type: 'POST',
              success: function(data) {
                if(data.status=="success"){
                    $.jGrowl("Successfull");

                }
              }
          });
      });


  });


  function updatelist(id) {

      if (id == "usertime") {
          var url = "/getmessages";
      } else {
          var url = "/getgroupdatabyid/" + id;
      }
      console.log(url);
      $.get(url, function(data) {
          sap = JSON.stringify(data);
          var length = data.data.length;
          console.log("from update list" + length);
          $("#usertimeline").empty();
          for (var i = 0; i < length; i++) {
              if (data.data[i].content_type == "text") {
                      var timeline = "<li id=\"li" + data.data[i].timestamp + "\"><div class=\"timeline-badge warning\"><i class=\"glyphicon glyphicon-envelope\"></i></div><div class=\"timeline-panel\">  <div id=" + data.data[i].timestamp + " class=\"close\">&times;</div>" + "<div class=\"timeline-heading\"><h4 class=\"timeline-title\">" + linkify(data.data[i].content) + "</h4><p><small class=\"text-muted\">" + "<i class=\"glyphicon glyphicon-time\"></i>" + timeConverter(data.data[i].timestamp) + "</small></p></div></div></li>"

                  } else if (data.data[i].content_type == "url") {
                      var timeline = "<li id=\"li" + data.data[i].timestamp + "\"><div class=\"timeline-badge info\"><i class=\"glyphicon glyphicon-link\"></i></div><div class=\"timeline-panel\"> <div id=" + data.data[i].timestamp + " class=\"close\">&times;</div>" + "<div class=\"timeline-heading\"><h4 class=\"timeline-title\">" + linkify(data.data[i].content) + "</h4><p><small class=\"text-muted\">" + "<i class=\"glyphicon glyphicon-time\"></i>" + timeConverter(data.data[i].timestamp) + "</small></p></div></div></li>"
                  } else if (data.data[i].content_type == "file") {

                      var timeline = "<li id=\"li" + data.data[i].timestamp + "\"><div class=\"timeline-badge danger\"><i class=\"glyphicon glyphicon-file\"></i></div><div class=\"timeline-panel\">   <div id=" + data.data[i].timestamp + " class=\"close\">&times;</div>" + "<div class=\"timeline-heading\"><h4 class=\"timeline-title\"><a href='https://s3-us-west-2.amazonaws.com/megh-uploads/" + data.data[i].content + "' class=\"glyphicon glyphicon-download-alt\">  </a> " + data.data[i].content + "</h4><p><small class=\"text-muted\">" + "<i class=\"glyphicon glyphicon-time\"></i>" + timeConverter(data.data[i].timestamp) + "</small></p></div></div></li>"

                  }
                  $("#usertimeline").append(timeline);
          }

      });
  }

  function sendFileToServer(formData, status, url, ug) {
      var uploadURL = url; //Upload URL
      var extraData = {}; //Extra Data.
      var jqXHR = $.ajax({
          xhr: function() {
              var xhrobj = $.ajaxSettings.xhr();
              if (xhrobj.upload) {
                  xhrobj.upload.addEventListener('progress', function(event) {
                      var percent = 0;
                      var position = event.loaded || event.position;
                      var total = event.total;
                      if (event.lengthComputable) {
                          percent = Math.ceil(position / total * 100);
                      }
                      //Set progress
                      status.setProgress(percent);
                  }, false);
              }
              return xhrobj;
          },
          url: uploadURL,
          type: "POST",
          contentType: false,
          processData: false,
          cache: false,
          data: formData,
          success: function(data) {
              status.setProgress(100);
              if (ug == "group") {
                  updatelist(currentgroupid);
                  chat.ws.send('{"type":"message","from":"sanath","groupid":"' + currentgroupid + '","text":"file upload"}');
              } else if (ug == "user") {
                  updatelist("usertime");
              }
          }
      });

      status.setAbort(jqXHR);
  }

  var rowCount = 0;

  function createStatusbar(obj) {
      rowCount++;
      var row = "odd";
      if (rowCount % 2 == 0) row = "even";
      this.statusbar = $("<div class='statusbar " + row + "'></div>");
      this.filename = $("<div class='filename'></div>").appendTo(this.statusbar);
      this.size = $("<div class='filesize'></div>").appendTo(this.statusbar);
      this.progressBar = $("<div class='progressBar'><div></div></div>").appendTo(this.statusbar);
      this.abort = $("<div class='abort'>Abort</div>").appendTo(this.statusbar);
      obj.after(this.statusbar);

      this.setFileNameSize = function(name, size) {
          var sizeStr = "";
          var sizeKB = size / 1024;
          if (parseInt(sizeKB) > 1024) {
              var sizeMB = sizeKB / 1024;
              sizeStr = sizeMB.toFixed(2) + " MB";
          } else {
              sizeStr = sizeKB.toFixed(2) + " KB";
          }

          this.filename.html(name);
          this.size.html(sizeStr);
      }
      this.setProgress = function(progress) {
          var progressBarWidth = progress * this.progressBar.width() / 100;
          this.progressBar.find('div').animate({
              width: progressBarWidth
          }, 10).html(progress + "% ");
          if (parseInt(progress) >= 100) {
              this.abort.hide();
          }
      }
      this.setAbort = function(jqxhr) {
          var sb = this.statusbar;
          this.abort.click(function() {
              jqxhr.abort();
              sb.hide();
          });
      }
  }

  function handleFileUpload(files, obj, url, ug) {
      for (var i = 0; i < files.length; i++) {
          var fd = new FormData();
          fd.append('file', files[i]);

          var status = new createStatusbar(obj); //Using this we can set progress.
          status.setFileNameSize(files[i].name, files[i].size);
          sendFileToServer(fd, status, url, ug);

      }
  }
  $(document).ready(function() {
      var obj = $("#userdragandrophandler");

      obj.on('dragenter', function(e) {
          e.stopPropagation();
          e.preventDefault();
          $(this).css('border', '2px solid #0B85A1');
      });


      obj.on('dragover', function(e) {
          e.stopPropagation();
          e.preventDefault();
      });

      obj.on('drop', function(e) {

          $(this).css('border', '2px dotted #0B85A1');
          e.preventDefault();
          var files = e.originalEvent.dataTransfer.files;
          if (currentgroupid == 0) {
              //We need to send dropped files to Server
              handleFileUpload(files, obj, "/uploadfile", "user");
          } else {
              //We need to send dropped files to Server
              handleFileUpload(files, obj, "/groupuploadfile/" + currentgroupid, "group");
          }

      });


      $(document).on('dragenter', function(e) {
          e.stopPropagation();
          e.preventDefault();
      });
      $(document).on('dragover', function(e) {
          e.stopPropagation();
          e.preventDefault();
          obj.css('border', '2px dotted #0B85A1');
      });
      $(document).on('drop', function(e) {
          e.stopPropagation();
          e.preventDefault();
      });

  });

  function linkify(text) {
      var urlRegex = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
      return text.replace(urlRegex, function(url) {
          return '<a href="' + url + '">' + url + '</a>';
      })
  }

  function checkurl(text) {
      var temp = text;
      var temp2 = text;
      var match = (text.match(/(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig) || []).length;
      console.log(match + "," + temp.length);
      if (match == 1) {
          var array = temp.match(/(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig);
          if (array[0].length == temp2.length) {
              return true;

          } else {
              return false;
          }
      } else {
          return false;
      }
  }

  function getCookies() {
      var c = document.cookie,
          v = 0,
          cookies = {};
      if (document.cookie.match(/^\s*\$Version=(?:"1"|1);\s*(.*)/)) {
          c = RegExp.$1;
          v = 1;
      }
      if (v === 0) {
          c.split(/[,;]/).map(function(cookie) {
              var parts = cookie.split(/=/, 2),
                  name = decodeURIComponent(parts[0].trimLeft()),
                  value = parts.length > 1 ? decodeURIComponent(parts[1].trimRight()) : null;
              cookies[name] = value;
          });
      } else {
          c.match(/(?:^|\s+)([!#$%&'*+\-.0-9A-Z^`a-z|~]+)=([!#$%&'*+\-.0-9A-Z^`a-z|~]*|"(?:[\x20-\x7E\x80\xFF]|\\[\x00-\x7F])*")(?=\s*[,;]|$)/g).map(function($0, $1) {
              var name = $0,
                  value = $1.charAt(0) === '"' ? $1.substr(1, -1).replace(/\\(.)/g, "$1") : $1;
              cookies[name] = value;
          });
      }
      return cookies;
  }

  function getCookie(name) {
      return getCookies()[name];
  }

  function timeConverter(UNIX_timestamp) {
      var a = new Date(UNIX_timestamp.slice(0, 10) * 1000);

      var months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
      var year = a.getFullYear();
      var month = months[a.getMonth()];
      var date = a.getDate();
      var hour = a.getHours();
      var min = a.getMinutes();
      var sec = a.getSeconds();
      var time = '  ' + hour + ':' + min + ':' + sec + '  , ' + date + '  ' + month + '   ' + year + ' ';
      return time;
  }
