package com.example.salman.login;

import org.json.JSONObject;

import java.util.ArrayList;
import java.util.List;

/**
 * Created by salman on 11/24/14.
 */
public class JsonCustomParser {
    public static List<Groups> readAndParseJSONGroups(String content) {
        try {
            JSONObject reader = new JSONObject(content);


            List<Groups> groupsList = new ArrayList<>();

            for (int i = 0; i < reader.getJSONArray("data").length(); i++){
                JSONObject obj = reader.getJSONArray("data").getJSONObject(i);
                Groups group = new Groups();

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

    public static List<Message> readAndParseJSONMessages(String content) {
        try {
            JSONObject reader = new JSONObject(content);


            List<Message> messageList = new ArrayList<>();

            for (int i = 0; i < reader.getJSONArray("data").length(); i++){
                JSONObject obj = reader.getJSONArray("data").getJSONObject(i);
                Message message = new Message();

                message.setContent(obj.getString("content"));
                message.setContent_type(obj.getString("content_type"));
                message.setTimestamp(obj.getString("timestamp"));

                messageList.add(message);
            }

            return messageList;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }
}