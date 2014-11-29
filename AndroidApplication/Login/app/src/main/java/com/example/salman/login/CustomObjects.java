package com.example.salman.login;

/**
 * Created by salman on 11/28/14.
 */
public class CustomObjects {

    static class Groups {
        private String group_id;
        private String username;
        private String group_name;
        private String group_admin;
        private String timestamp;


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
    }

    static class GroupData {

        private String groupdata_id;
        private String content;
        private String content_type;
        private String username;
        private String timestamp;

        public String getGroupdata_id() {
            return groupdata_id;
        }

        public void setGroupdata_id(String groupdata_id) {
            this.groupdata_id = groupdata_id;
        }

        public String getContent() {
            return content;
        }

        public void setContent(String content) {
            this.content = content;
        }

        public String getUsername() {
            return username;
        }

        public void setUsername(String username) {
            this.username = username;
        }

        public String getTimestamp() {
            return timestamp;
        }

        public void setTimestamp(String timestamp) {
            this.timestamp = timestamp;
        }

        public String getContent_type() {
            return content_type;
        }

        public void setContent_type(String content_type) {
            this.content_type = content_type;
        }
    }

    static class UserMessage {
        private String content, content_type, timestamp;

        public String getContent() {
            return content;
        }

        public void setContent(String content) {
            this.content = content;
        }

        public String getContent_type() {
            return content_type;
        }

        public void setContent_type(String content_type) {
            this.content_type = content_type;
        }

        public String getTimestamp() {
            return timestamp;
        }

        public void setTimestamp(String timestamp) {
            this.timestamp = timestamp;
        }
    }

    static class SelectedContent {
        private String content;
        private int position;

        public int getPosition() {
            return position;
        }

        public void setPosition(int position) {
            this.position = position;
        }

        public String getContent() {
            return content;
        }

        public void setContent(String content) {
            this.content = content;
        }
    }

}
