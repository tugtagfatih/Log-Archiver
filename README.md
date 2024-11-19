# **_tugtagfatih's_ DevOps Tools**

## [Log Archiver](https://roadmap.sh/projects/log-archive-tool)
- Provide the log directory as an argument when running the tool. <br /> ``` log-archive log-directory> ``` <br /> If the "log-directory" variable is not entered, the ```/var/log``` directory will be selected by default.

- The archived version of your files is saved in the ```log-directory/archived_logs``` as ```logs_archive_date_time.tar.gz ``` <br />  **Example:** ```📁logs_archive_20241119_121631.tar.gz```

> [!IMPORTANT]
> If you are using the default log directory or the directory where your logs are located requires root privilege, run the program with root privilege.

