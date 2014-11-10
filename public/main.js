  var currentgroupid;
    var currentgroupname;

 var urlRegex =/(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;
  var ChatRoom = function() {
    var c = this;
    
    c.create = function() {
      c.sendForm   = document.getElementById("send-form");
      c.textInput  = c.sendForm.querySelector("#groupmessage");
      
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
        console.log("received"+e.data);
        
        var temp=JSON.parse(e.data);

        if(currentgroupid==temp.groupid)
        {
          updatelist(temp.groupid);
        }
        else{
          val=$("#badge"+temp.groupid).text();
          if(val==undefined)
          {
            val=0;
          }
          val=val+1;
          $("#badge"+temp.groupid).text(val);
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
      console.log("sent"+JSON.stringify(messageData));
      c.ws.send(JSON.stringify(messageData));
    };

    c.handleSend = function() {
      c.sendForm.addEventListener("click", function(e) {
        e.preventDefault();
        
        if (c.textInput.value && c.textInput.value.length) {
          var messageData = {
            type: "message",
            from : "sanath",
            groupid:""+currentgroupid,
            text: c.textInput.value
          };        
          c.send(messageData);
          c.textInput.value = ''
        }      
        return false;
      });    
      return c;
    };

    var initialize = function() {
      c.create().handleSend();
      c.connect().listen();
    }();
  };

  $( document ).ready(function() {
    chat = new ChatRoom();

    $('#publickey').text("publickey : "+getCookie("publickey"))
    $.get( "/getmessages", function( data ) {
      sap = JSON.stringify(data);
      var length = data.data.length;
      console.log("from userdata"+length);
      for(var i=0;i<length;i++)
      {
        if(data.data[i].content_type=="text"){
          var temp="<tr><td>"+data.data[i].content+" <td> <td align='right'>message</td></tr>";
        }
        else if(data.data[i].content_type=="file") {
          var temp="<tr><td>"+data.data[i].content+"<td> <td align='right'><a href='https://s3-us-west-2.amazonaws.com/megh-uploads/"+data.data[i].content+"'>download</a></td></tr>";
        }
        $("#user").append(temp);
      }

      $( "#result" ).append(sap);
    });
    $(document).on('click', '.groupclick',function(){
      var id =this.id;
      $("#badge"+id).text("");
      $(this).css("background-color", "#DCF6F7");
      if(currentgroupid!=id)
      {
        $("#"+currentgroupid+".groupclick").css("background-color", "white");
      }
      currentgroupid=id;
      currentgroupname=$(this).text();
      var url = "/getgroupdatabyid/"+id;
      $.get( url, function( data ) {
        sap = JSON.stringify(data);
        var length = data.data.length;
        console.log("from getgroupdata "+length);
        $("#group").empty();
        for(var i=0;i<length;i++)
        {
          if(data.data[i].content_type=="text"){
           var temp="<tr><td>"+data.data[i].content+" <td> <td align='right'>message</td></tr>";
         }
         else if(data.data[i].content_type=="file") {
          var temp="<tr><td>"+data.data[i].content+"<td> <td align='right'><a href='https://s3-us-west-2.amazonaws.com/megh-uploads/"+data.data[i].content+"'>download</a></td></tr>";
        }
        $("#group").append(temp);
      }

    });
    });

    $.get( "/groupslist", function( data ) {
      sap = JSON.stringify(data);
      var length = data.data.length;

      for(var i=0;i<length;i++)
      {
        if(i==0)
        {
          currentgroupid=data.data[i].group_id;
          currentgroupname=data.data[i].group_name;
          updatelist(currentgroupid);
        }
        var temp ="<li class=\"list-group-item groupclick\" id="+data.data[i].group_id+" align=\"center\"> <span id=\"badge"+data.data[i].group_id+"\" class=\"badge\"></span> "+data.data[i].group_name+"</li>";
      //var temp = "<option value="+data.data[i].groupid+">"+data.data[i].groupname+"</option>";
      $("#addgroups").append(temp);

      $("#"+currentgroupid+".groupclick").css("background-color", "#DCF6F7");
    }

  });

    $("#groupsubmit").click(function(){
      var value =$("#groupmessage").val();
      var id =currentgroupid;
      var url = "/groupdata/"+id;
      $.post(url,{message:value},function(data,status){
        //location.reload();
        $("#groupmessage").val("");
        console.log("submitted  "+status);
        updatelist(id);
        
      });
      
    });

    $("#creategroupsubmit").click(function(){
      var value =$("#groupname").val();

      var url = "/creategroup"
      $.post(url,{message:value},function(data,status){
        //location.reload();
        console.log(status);
        $("#groupname").val("");
        updategrouplist();
      });


    });


$("#addusers").click(function(){
    $("#addusers").toggle();
     $("#adduserbox").toggle();
});
    $("#adduser").click(function(){
      var addusername =$("#addusername").val();
      var addpublickey =$("#addpublickey").val();

      var url = "/adduser"
      $.post(url,{email:addusername,publickey:addpublickey,groupname:currentgroupname,groupid:currentgroupid},function(data,status){
        //location.reload();
        console.log(status);
        $("#addusers").toggle();
     $("#adduserbox").toggle();
        
      });


    });

    $("#userfilesend").click(function(){
      var value =$("#text-input").val();

      var url = "/messages"
      $.post(url,{message:value},function(data,status){
        //location.reload();
        console.log(status);
        $("#text-input").val("");
        updateuserlist();
      });
    });


    $("#groupfile").change(function(){
      var value =$("#groupmessage").val();
      var id =currentgroupid;
      var url = "/groupuploadfile/"+id;
      var file = $("#groupfile")[0].files[0];
      var formdata = new FormData();
      formdata.append("file", file);
      $.ajax({
        url: url,
        data: formdata,
        processData: false,
        contentType: false,
        type: 'POST',
        success: function(data) {
          updatelist(id);
          chat.ws.send('{"type":"message","from":"sanath","groupid":"'+currentgroupid+'","text":"file upload"}' );
        }
      });
    });

    $("#userfile").change(function(){
      var url = "/uploadfile"
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
          $("#uploadform").toggle();
          updateuserlist();
        }
      });
    });   

    $("#filebutton").click(function(){
      $("#uploadform").toggle();
    }); 

    $("#groupfilebutton").click(function(){
      $("#uploadformgroup").toggle();
    }); 

  });


  function updatelist(id){
    currentgroupid=id;
    var url = "/getgroupdatabyid/"+id;
    $.get( url, function( data ) {
      sap = JSON.stringify(data);
      var length = data.data.length;
      console.log("from update list"+length);
      $("#group").empty();
      for(var i=0;i<length;i++)
      {
        if(data.data[i].content_type=="text"){
         var temp="<tr><td>"+data.data[i].content+" <td> <td align='right'>message</td></tr>";
       }
       else if(data.data[i].content_type=="file") {
        var temp="<tr><td>"+data.data[i].content+"<td> <td align='right'><a href='https://s3-us-west-2.amazonaws.com/megh-uploads/"+data.data[i].content+"'>download</a></td></tr>";
      }
      $("#group").append(temp);
    }
    
  });
  }

  function updateuserlist(){
    var url = "/getmessages"
    $.get( url, function( data ) {
      sap = JSON.stringify(data);
      var length = data.data.length;
      console.log("from user list"+length);
      $("#user").empty();
      for(var i=0;i<length;i++)
      {
        if(data.data[i].content_type=="text"){
         var temp="<tr><td>"+data.data[i].content+" <td> <td align='right'>message</td></tr>";
       }
       else if(data.data[i].content_type=="file") {
        var temp="<tr><td>"+data.data[i].content+"<td> <td align='right'><a href='https://s3-us-west-2.amazonaws.com/megh-uploads/"+data.data[i].content+"'>download</a></td></tr>";
      }
      $("#user").append(temp);
    }
    
  });
  }

  function updategrouplist(){
    $.get( "/groupslist", function( data ) {
      sap = JSON.stringify(data);
      var length = data.data.length;
      var temp ="<li class=\"list-group-item groupclick\" id="+data.data[length-1].group_id+" align=\"center\"><span id=\"badge"+data.data[length-1].group_id+"\" class=\"badge\"></span> "+data.data[length-1].group_name+"</li>";
      //var temp = "<option value="+data.data[i].groupid+">"+data.data[i].groupname+"</option>";
      $("#addgroups").append(temp);
    });
  }

  function sendFileToServer(formData,status,url,ug)
  {
      var uploadURL =url; //Upload URL
      var extraData ={}; //Extra Data.
      var jqXHR=$.ajax({
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
        contentType:false,
        processData: false,
        cache: false,
        data: formData,
        success: function(data){
          status.setProgress(100);
          if(ug=="group"){
            updatelist(currentgroupid);
            chat.ws.send('{"type":"message","from":"sanath","groupid":"'+currentgroupid+'","text":"file upload"}' );
          }
          else if(ug=="user"){
            updateuserlist();                
          }
        }
      }); 

  status.setAbort(jqXHR);
}

var rowCount=0;
function createStatusbar(obj)
{
 rowCount++;
 var row="odd";
 if(rowCount %2 ==0) row ="even";
 this.statusbar = $("<div class='statusbar "+row+"'></div>");
 this.filename = $("<div class='filename'></div>").appendTo(this.statusbar);
 this.size = $("<div class='filesize'></div>").appendTo(this.statusbar);
 this.progressBar = $("<div class='progressBar'><div></div></div>").appendTo(this.statusbar);
 this.abort = $("<div class='abort'>Abort</div>").appendTo(this.statusbar);
 obj.after(this.statusbar);

 this.setFileNameSize = function(name,size)
 {
  var sizeStr="";
  var sizeKB = size/1024;
  if(parseInt(sizeKB) > 1024)
  {
    var sizeMB = sizeKB/1024;
    sizeStr = sizeMB.toFixed(2)+" MB";
  }
  else
  {
    sizeStr = sizeKB.toFixed(2)+" KB";
  }

  this.filename.html(name);
  this.size.html(sizeStr);
}
this.setProgress = function(progress)
{       
  var progressBarWidth =progress*this.progressBar.width()/ 100;  
  this.progressBar.find('div').animate({ width: progressBarWidth }, 10).html(progress + "% ");
  if(parseInt(progress) >= 100)
  {
    this.abort.hide();
  }
}
this.setAbort = function(jqxhr)
{
  var sb = this.statusbar;
  this.abort.click(function()
  {
    jqxhr.abort();
    sb.hide();
  });
}
}
function handleFileUpload(files,obj,url,ug)
{
 for (var i = 0; i < files.length; i++) 
 {
  var fd = new FormData();
  fd.append('file', files[i]);

          var status = new createStatusbar(obj); //Using this we can set progress.
          status.setFileNameSize(files[i].name,files[i].size);
          sendFileToServer(fd,status,url,ug);

        }
      }
      $(document).ready(function()
      {
        var obj = $("#userdragandrophandler");

        var gobj = $("#groupdragandrophandler");
        obj.on('dragenter', function (e) 
        {
          e.stopPropagation();
          e.preventDefault();
          $(this).css('border', '2px solid #0B85A1');
        });
        gobj.on('dragenter', function (e) 
        {
          e.stopPropagation();
          e.preventDefault();
          $(this).css('border', '2px solid #0B85A1');
        });

        obj.on('dragover', function (e) 
        {
         e.stopPropagation();
         e.preventDefault();
       });
        gobj.on('dragover', function (e) 
        {
         e.stopPropagation();
         e.preventDefault();
       });
        obj.on('drop', function (e) 
        {

         $(this).css('border', '2px dotted #0B85A1');
         e.preventDefault();
         var files = e.originalEvent.dataTransfer.files;

       //We need to send dropped files to Server
       handleFileUpload(files,obj,"/uploadfile","user");
     });
        gobj.on('drop', function (e) 
        {

         $(this).css('border', '2px dotted #0B85A1');
         e.preventDefault();
         var files = e.originalEvent.dataTransfer.files;

       //We need to send dropped files to Server
       handleFileUpload(files,gobj,"/groupuploadfile/"+currentgroupid,"group");
     });

        $(document).on('dragenter', function (e) 
        {
          e.stopPropagation();
          e.preventDefault();
        });
        $(document).on('dragover', function (e) 
        {
          e.stopPropagation();
          e.preventDefault();
          obj.css('border', '2px dotted #0B85A1');
        });
        $(document).on('drop', function (e) 
        {
          e.stopPropagation();
          e.preventDefault();
        });

      });

  function linkify(text) {  
        var urlRegex =/(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;  
        return text.replace(urlRegex, function(url) {  
                            return '<a href="' + url + '">' + url + '</a>';  
          })  
  }

  if (typeof String.prototype.trimLeft !== "function") {
    String.prototype.trimLeft = function() {
        return this.replace(/^\s+/, "");
    };
}
if (typeof String.prototype.trimRight !== "function") {
    String.prototype.trimRight = function() {
        return this.replace(/\s+$/, "");
    };
}
if (typeof Array.prototype.map !== "function") {
    Array.prototype.map = function(callback, thisArg) {
        for (var i=0, n=this.length, a=[]; i<n; i++) {
            if (i in this) a[i] = callback.call(thisArg, this[i]);
        }
        return a;
    };
}
function getCookies() {
    var c = document.cookie, v = 0, cookies = {};
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
                value = $1.charAt(0) === '"'
                          ? $1.substr(1, -1).replace(/\\(.)/g, "$1")
                          : $1;
            cookies[name] = value;
        });
    }
    return cookies;
}
function getCookie(name) {
    return getCookies()[name];
}