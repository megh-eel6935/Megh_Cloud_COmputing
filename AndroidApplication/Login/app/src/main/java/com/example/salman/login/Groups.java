package com.example.salman.login;

/**
 * Created by salman on 11/27/14.
 */
public class Groups {
    private String group_id;
    private String username;
    private String group_name;
    private String group_admin;

    public String getTimestamp() {
        return timestamp;
    }

    public void setTimestamp(String timestamp) {
        this.timestamp = timestamp;
    }

    public String getGroup_admin() {
        return group_admin;
    }

    public void setGroup_admin(String group_admin) {
        this.group_admin = group_admin;
    }

    public String getGroup_name() {
        return group_name;
    }

    public void setGroup_name(String group_name) {
        this.group_name = group_name;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getGroup_id() {
        return group_id;
    }

    public void setGroup_id(String group_id) {
        this.group_id = group_id;
    }

    private String timestamp;
}

