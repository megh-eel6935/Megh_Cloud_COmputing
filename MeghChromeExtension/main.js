chrome.contextMenus.removeAll();

function sendurlOnClick(info, tab) {
    console.log("item " + info.menuItemId + " was clicked");
    console.log("info: " + JSON.stringify(info));
    console.log("tab: " + JSON.stringify(tab));

    chrome.tabs.query({
        currentWindow: true,
        active: true
    }, function(tabs) {
       
       $.post("http://104.131.126.89/urls", {
            message: tabs[0].url

        }, function(data, status) {
        	console.log(data);
            if (data.status == "success") {

                alert("url sent successfully");

            }
        });
    });
}

function sendtextOnClick(info, tab) {
    console.log("item " + info.menuItemId + " was clicked");
    console.log("info: " + JSON.stringify(info));
    console.log("tab: " + JSON.stringify(tab));


    $.post("http://104.131.126.89/messages", {
        message: info.selectionText

    }, function(data, status) {
        if (data.status == "success") {

            alert("text sent successfully");

        }
    });


}

function sendfileOnClick(info, tab) {
    console.log("item " + info.menuItemId + " was clicked");
    console.log("info: " + JSON.stringify(info));
    console.log("tab: " + JSON.stringify(tab));
    var str = info.srcUrl;
    var type ="somefile"
    if(str.indexOf('.jpg') !== -1){
    	type=".jpg";
    }
    else if(str.indexOf('.png') !== -1){
    	type=".png";
    }
    else if(str.indexOf('.pdf') !== -1){
    	type=".pdf";
    }
    else if(str.indexOf('.txt') !== -1){
    	type=".txt";
    }
    else if(str.indexOf('.mp3') !== -1){
    	type=".mp3";
    }
    else if(str.indexOf('.mp4') !== -1){
    	type=".mp4";
    }

    $.post("http://104.131.126.89/uploadfilefromurl", {
        url: info.srcUrl,
        type: type
    }, function(data, status) {
        console.log(status);
        if (data.status == "success") {

            alert("file sent successfully");

        }
    });

}

// Create one test item for each context type.

function groupurlOnClick(info, tab) {

    console.log("item " + info.menuItemId + " was clicked");
    console.log("info: " + JSON.stringify(info));
    console.log("tab: " + JSON.stringify(tab));

    chrome.tabs.query({
        currentWindow: true,
        active: true
    }, function(tabs) {
        //alert("sendurl"+tabs[0].url);  

        var groupid = info.menuItemId.substr(info.menuItemId.indexOf(";") + 1);
        $.post("http://104.131.126.89/groupdataurl/" + groupid, {
            message: tabs[0].url

        }, function(data, status) {
            if (data.status == "success") {

                alert("url sent successfully");

            }
        });
    });
}

function grouptextOnClick(info, tab) {

    console.log("item " + info.menuItemId + " was clicked");
    console.log("info: " + JSON.stringify(info));
    console.log("tab: " + JSON.stringify(tab));

    var groupid = info.menuItemId.substr(info.menuItemId.indexOf(";") + 1);

    $.post("http://104.131.126.89/groupdata/" + groupid, {
        message: info.selectionText

    }, function(data, status) {
        if (data.status == "success") {

            alert("text sent successfully");

        }
    });
}



function groupFileOnClick(info, tab) {

    console.log("item " + info.menuItemId + " was clicked");
    console.log("info: " + JSON.stringify(info));
    console.log("tab: " + JSON.stringify(tab));

    var groupid = info.menuItemId.substr(info.menuItemId.indexOf(";") + 1);

    var str = info.srcUrl;
    var type ="somefile"
    if(str.indexOf('.jpg') !== -1){
    	type=".jpg";
    }
    else if(str.indexOf('.png') !== -1){
    	type=".png";
    }
    else if(str.indexOf('.pdf') !== -1){
    	type=".pdf";
    }
    else if(str.indexOf('.txt') !== -1){
    	type=".txt";
    }
    else if(str.indexOf('.mp3') !== -1){
    	type=".mp3";
    }
    else if(str.indexOf('.mp4') !== -1){
    	type=".mp4";
    }

    $.post("http://104.131.126.89/groupuploadfilefromurl/" + groupid, {
        url: info.srcUrl,
        type: type
    }, function(data, status) {
        console.log(status);
        if (data.status == "success") {

            alert("file sent successfully");

        }
    });
}


var childs = [];

$.get("http://104.131.126.89/groupslist", function(data) {
    sap = JSON.stringify(data);
    console.log(data);

    if (data.status == "Access denied") {

    	var sendurl = chrome.contextMenus.create({
            "title": "Login to get options",
            contexts: ["all"]
        });

    } else {

        var sendurl = chrome.contextMenus.create({
            "title": "sendurl",
            contexts: ["all"],
            "onclick": sendurlOnClick
        });

        var sendtext = chrome.contextMenus.create({
            "title": "sendtext",
            contexts: ["selection"],
            "onclick": sendtextOnClick
        });

        var sendfile = chrome.contextMenus.create({
            "title": "sendfile",
            contexts: ["image", "video", "audio"],
            "onclick": sendfileOnClick
        });

        var length = data.data.length;

        for (var i = 0; i < length; i++) {

            //if (data.data[i].username == data.data[i].group_admin) {
            // var temp = "<option id=" + data.data[i].group_id + " value=" + data.data[i].group_id + ">" + data.data[i].group_name + "</option>";
            var name = data.data[i].group_name;
            var group_id = data.data[i].group_id;
            var id = chrome.contextMenus.create({
                "title": "Group: " + name,
                contexts: ["all"]
            });

            var id1 = chrome.contextMenus.create({
                "id": "url;" + group_id,
                "title": "sendurl",
                contexts: ["all"],
                "parentId": id,
                "onclick": groupurlOnClick
            });

            var id2 = chrome.contextMenus.create({
                "id": "text;" + group_id,
                "title": "sendtext",
                contexts: ["selection"],
                "parentId": id,
                "onclick": grouptextOnClick
            });


            var id3 = chrome.contextMenus.create({
                "id": "file;" + group_id,
                "title": "send file",
                contexts: ["image", "video", "audio"],
                "parentId": id,
                "onclick": groupFileOnClick
            });

        }
    }

});
