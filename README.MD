# Project Proposal
by Jeremiah Cheng, Andrew Chan, Anna Zhou, Jimmy Hua
# Project Description
This application is intended to be used by people who need to easily figure out when other people in a group are free. This could be used for scheduling meetings, but also for any sort of group activity. Whether it be grabbing coffee, having a party or planning a company meeting, this app will provide the capability for people to easily find out when every person is free. Realistically, this will be used most by students and young adults, and potentially working professionals for planning activities. This application is useful because it allows for users to log in and save their schedule information and allows for the planning of multiple events. This has more functionality than current alternatives such as When2Meet which only has temporary meetings. It should make planning events and get-togethers even easier and much more efficient. We as developers would like to build this application as it will give us more practice in authenticating users and allowing for them to interact with other users. It should present a valuable experience in furthering our knowledge of Golang and its capabilities alongside Docker and associated Databases.

# Technical Description
![](imgs/diagram.png)
![](imgs/priority.png)

### User Fields:
User ID,
User Name,
Meeting IDs (array),
Free Times (array)

### Meeting Fields:
Meeting ID,
Meeting Name,
User IDs (array),
Free Times (array), calculated from users’ free times


The schedule will be constructed based on the free times
![](imgs/endpoint.png)
