package com.example.salman.login;

import org.json.JSONObject;

import java.util.ArrayList;
import java.util.List;

/**
 * Created by salman on 11/24/14.
 */
public class JsonCustomParser {
    public static List<CustomObjects.Groups> readAndParseJSONGroups(String content) {
        try {
            JSONObject reader = new JSONObject(content);


            List<CustomObjects.Groups> groupsList = new ArrayList<>();

            for (int i = 0; i < reader.getJSONArray("data").length(); i++){
                JSONObject obj = reader.getJSONArray("data").getJSONObject(i);
                CustomObjects.Groups group = new CustomObjects.Groups();

                group.setGroup_id(obj.getString("group_id"));
                group.setGroup_admin(obj.getString("group_admin"));
                group.setGroup_name(obj.getString("group_name"));
                group.setUsername(obj.getString("username"));
                group.setTimestamp(obj.getString("timestamp"));

                groupsList.add(group);
            }

            return groupsList;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    public static List<CustomObjects.GroupData> readAndParseJSONGroupData(String content) {
        try {
            JSONObject reader = new JSONObject(content);


            List<CustomObjects.GroupData> groupDataList = new ArrayList<>();

            for (int i = 0; i < reader.getJSONArray("data").length(); i++){
                JSONObject obj = reader.getJSONArray("data").getJSONObject(i);
                CustomObjects.GroupData groupData = new CustomObjects.GroupData();

                groupData.setGroupdata_id(obj.getString("groupdata_id"));
                groupData.setContent(obj.getString("content"));
                groupData.setContent_type(obj.getString("content_type"));
                groupData.setUsername(obj.getString("username"));
                groupData.setTimestamp(obj.getString("timestamp"));

                groupDataList.add(groupData);
            }

            return groupDataList;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    public static List<CustomObjects.UserMessage> readAndParseJSONMessages(String content) {
        try {
            JSONObject reader = new JSONObject(content);


            List<CustomObjects.UserMessage> userMessageList = new ArrayList<>();

            for (int i = 0; i < reader.getJSONArray("data").length(); i++){
                JSONObject obj = reader.getJSONArray("data").getJSONObject(i);
                CustomObjects.UserMessage userMessage = new CustomObjects.UserMessage();

                userMessage.setContent(obj.getString("content"));
                userMessage.setContent_type(obj.getString("content_type"));
                userMessage.setTimestamp(obj.getString("timestamp"));

                userMessageList.add(userMessage);
            }

            return userMessageList;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }
}