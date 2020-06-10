# Extermal interfaces

* Chatbot on AWS
    * Interface: HTTP/JSON
* Jackett search
    * Interface: HTTP/XML
* Transmission
    * Interface: HTTP/JSON
* Datastorage
    * Interface: TBA
* Kodi
    * Interface: indirect/SMB fileshare

# Action flow

1. Server receives search request from Chatbot
2. Server sends a search request to Jackett
3. Server receives a response with search results from Jackett
4. Server records the responses to Datastorage
5. Server forms response based on search results by discarding unnecessary fields
5. Server sends response to Chatbot
6. Server receives download requests from Chatbot
7. Server fetches search result matching the download request from Datastorage
8. Server sends download command to Transmission
9. Server records the download process to Datastorage
10. Server starts a process to peridiocally check the download progress from Transmission
10. Server sends response to Chatbot with the download process ID
11. Server receives download status query from Chatbot
12. Server fetches the download process information from Datastorage
14. Server sends response to Chatbot telling the current download progress
15. Once the process checking the download status notices the download is complete it stops the download
16. Server updates datastorage that download is complete
17. Server names the downloaded files so they can be scraped by Kodi
17. Server moves the downloaded files to the file server location
18. Server receives another download status query from Chatbot
19. Server fetches the download process information from Datastorage
20. Server sends response to Chatbot telling the download is complete