package com.example.salman.login;

import android.app.Activity;
import android.content.Context;
import android.content.res.Configuration;
import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v4.widget.DrawerLayout;
import android.support.v7.app.ActionBarActivity;
import android.support.v7.app.ActionBarDrawerToggle;
import android.view.LayoutInflater;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.BaseAdapter;
import android.widget.ImageView;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

import java.util.List;

/**
 * Created by salman on 11/27/14.
 */
public class DrawerActivity extends ActionBarActivity {

    private DrawerLayout drawerLayout;
    private ListView drawerLeftListView;
    private ListView drawerRightListView;

    private ActionBarDrawerToggle drawerLeftListener;

    private AdapterLeftDrawerList adapterLeftDrawerList;
    private AdapterMessage adapterMessage;
    private AdapterGroups adapterGroups;

    private TextView testTextView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_drawer);

        drawerLayout = (DrawerLayout) findViewById(R.id.drawer_activityID);
        drawerLeftListView = (ListView) findViewById(R.id.drawer_left_slider_listView);
        drawerRightListView = (ListView) findViewById(R.id.drawer_right_slider_listView);

        //////////////////
        testTextView = (TextView) findViewById(R.id.drawer_insideTextView);
        /////////////////

        adapterLeftDrawerList = new AdapterLeftDrawerList(this);
        drawerLeftListView.setAdapter(adapterLeftDrawerList);

        drawerLeftListView.setOnItemClickListener(new ItemClickListener());

        drawerLeftListener = new ActionBarDrawerToggle(this, drawerLayout,
                R.string.drawer_open, R.string.drawer_close) {

            @Override
            public void onDrawerOpened(View drawerView) {
                //super.onDrawerOpened(drawerView);
                Toast.makeText(DrawerActivity.this, "Drawer Opened", Toast.LENGTH_SHORT).show();
            }

            @Override
            public void onDrawerClosed(View drawerView) {
                //super.onDrawerClosed(drawerView);
                Toast.makeText(DrawerActivity.this, "Drawer Closed", Toast.LENGTH_SHORT).show();
            }
        };

        drawerLayout.setDrawerListener(drawerLeftListener);
        getSupportActionBar().setHomeButtonEnabled(true);
        getSupportActionBar().setDisplayHomeAsUpEnabled(true);
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        if (drawerLeftListener.onOptionsItemSelected(item)) {
            return true;
        }
        return super.onOptionsItemSelected(item);
    }

    @Override
    public void onConfigurationChanged(Configuration newConfig) {
        super.onConfigurationChanged(newConfig);
        drawerLeftListener.onConfigurationChanged(newConfig);
    }

    @Override
    protected void onPostCreate(Bundle savedInstanceState) {
        super.onPostCreate(savedInstanceState);
        drawerLeftListener.syncState();
    }

    private class ItemClickListener implements AdapterView.OnItemClickListener {

        @Override
        public void onItemClick(AdapterView<?> parent, View view, int position, long id) {

            String content = null;

            Toast.makeText(DrawerActivity.this, adapterLeftDrawerList.leftDrawerOptions[position] + " was selected",
                    Toast.LENGTH_SHORT).show();

            BTGetSelectedItemFunctionalityHandler itemClicked = new BTGetSelectedItemFunctionalityHandler();
            itemClicked.execute(position);

            selectMethod(position);
            drawerLayout.closeDrawer(drawerLeftListView);
        }

        private class BTGetSelectedItemFunctionalityHandler extends AsyncTask<Integer, String, SelectedContent> {

            @Override
            protected SelectedContent doInBackground(Integer... params) {
                try {
                    String content = UserFunctions.getSelectedListItemFunctionality(params[0]);
                    int position = params[0];

                    SelectedContent selectedContent = new SelectedContent();
                    selectedContent.setContent(content);
                    selectedContent.setPosition(position);
                    return selectedContent;

                } catch (Exception e) {
                    e.printStackTrace();
                    return null;
                }
            }

            @Override
            protected void onPostExecute(SelectedContent s) {
                ListView drawer_baseLayout_listView = (ListView) findViewById(R.id.drawer_baseLayout_listView);
                if (s.getPosition() == 0) {
                    List<Message> messageList = JsonCustomParser.readAndParseJSONMessages(s.getContent());
                    adapterMessage = new AdapterMessage(DrawerActivity.this,
                            R.layout.internallayout_drawer_base_item, messageList);

                    drawer_baseLayout_listView.setAdapter(adapterMessage);

                } else if (s.getPosition() == 1) {
                    List<Message> messageList = JsonCustomParser.readAndParseJSONMessages(s.getContent());
                    adapterMessage = new AdapterMessage(DrawerActivity.this,
                            R.layout.internallayout_drawer_base_item, messageList);

                    drawer_baseLayout_listView.setAdapter(adapterMessage);

                } else if (s.getPosition() == 2){
                    List<Groups> groupsList = JsonCustomParser.readAndParseJSONGroups(s.getContent());
                    adapterGroups = new AdapterGroups(DrawerActivity.this,
                            R.layout.internallayout_drawer_base_item, groupsList);

                    drawer_baseLayout_listView.setAdapter(adapterGroups);
                }
            }
        }
    }

    public void selectMethod(int position) {
        drawerLeftListView.setItemChecked(position, true);
        setTitle(adapterLeftDrawerList.leftDrawerOptions[position]);
    }

    public void setTitle(String title) {
        getSupportActionBar().setTitle(title);
    }
}

class AdapterLeftDrawerList extends BaseAdapter {

    private Context context;
    String[] leftDrawerOptions;
    int[] images = {R.drawable.ic_home, R.drawable.ic_people, R.drawable.ic_groups};


    public AdapterLeftDrawerList(Context context) {
        this.context = context;
        leftDrawerOptions = context.getResources().getStringArray(R.array.nav_drawer_left_items);
    }

    @Override
    public int getCount() {
        return leftDrawerOptions.length;
    }

    @Override
    public Object getItem(int position) {
        return leftDrawerOptions[position];
    }

    @Override
    public long getItemId(int position) {
        return position;
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        View row;
        if (convertView == null) {
            LayoutInflater inflater = (LayoutInflater) context.getSystemService(Context.LAYOUT_INFLATER_SERVICE);
            row = inflater.inflate(R.layout.internallayout_drawer_listitem, parent, false);
        } else {
            row = convertView;
        }
        TextView titleTextView = (TextView) row.findViewById(R.id.internalLayout_listItem_title_textView);
        ImageView titleImageView = (ImageView) row.findViewById(R.id.internalLayout_listItem_icon_imageView);

        titleTextView.setText(leftDrawerOptions[position]);
        titleImageView.setImageResource(images[position]);

        return row;
    }
}


class AdapterMessage extends ArrayAdapter<Message> {

    private Context context;
    private List<Message> messageList;

    public AdapterMessage(Context context, int resource, List<Message> objects) {
        super(context, resource, objects);
        this.context = context;
        this.messageList = objects;
    }

    public View getView(int position, View convertView, ViewGroup parent) {
        LayoutInflater inflater =
                (LayoutInflater) context.getSystemService(Activity.LAYOUT_INFLATER_SERVICE);

        View view = inflater.inflate(R.layout.internallayout_drawer_base_item, parent, false);

        Message msg = messageList.get(position);
        TextView tv = (TextView) view.findViewById(R.id.internalLayout_baseItem_title_textView);
        tv.setText(msg.getContent());


        ImageView img = (ImageView) view.findViewById(R.id.internalLayout_baseItem_icon_imageView);

        if (msg.getContent_type().equals("text")) {
            img.setImageResource(R.drawable.ic_message);
        } else if (msg.getContent_type().equals("file")) {
            img.setImageResource(R.drawable.ic_file);
        } else if (msg.getContent_type().equals("url")) {
            img.setImageResource(R.drawable.ic_url);
        }

        return view;
    }
}

class AdapterGroups extends ArrayAdapter<Groups> {

    private Context context;
    private List<Groups> groupsList;

    public AdapterGroups(Context context, int resource, List<Groups> objects) {
        super(context, resource, objects);
        this.context = context;
        this.groupsList = objects;
    }

    public View getView(int position, View convertView, ViewGroup parent) {
        LayoutInflater inflater =
                (LayoutInflater) context.getSystemService(Activity.LAYOUT_INFLATER_SERVICE);

        View view = inflater.inflate(R.layout.internallayout_drawer_base_item, parent, false);

        Groups grp = groupsList.get(position);
        TextView tv = (TextView) view.findViewById(R.id.internalLayout_baseItem_title_textView);
        tv.setText(grp.getGroup_name());


        ImageView img = (ImageView) view.findViewById(R.id.internalLayout_baseItem_icon_imageView);

        img.setImageResource(R.drawable.ic_groups);

        return view;

    }
}

class SelectedContent {
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
