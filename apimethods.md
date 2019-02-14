# Chat Api Functionality

1) Create New User :

       ```
       mutation{
         newUser(input:{
           name:"xyz",
           email:"xyz@gmail.com",
           contact:"0123456789",
           profilePicture:"xyz"
           bio:"xyz"
         }){
           name
           email
           contact
           profilePicture
           bio
           createdAt
         }
       }
       ```
2) Retrieve All Users Details Except Own:

    ```
    query{
      users(name:"xyz"){
        id
        name
        email
        contact
        profilePicture
        bio
        createdAt
      }
    }
    ```
3) Get Live Update For New User Join:

    ```
    subscription{
      userJoined{
        id
        name
        email
        contact
        profilePicture
        bio
        createdAt
      }
    }
    ``` 
4) Create New ChatRoom:
   
    ```
    mutation{
      newChatRoom(input:{
        creatorID:"84851515",
        chatRoomType:PRIVATE/GROUP
      },receiver:"48485454"){
        creatorID
        chatRoomID
        creator{
          name
          email
          contact
          profilePicture
          bio
          createdAt
        }
        chatRoomName
        chatRoomType
        createdAt
      }
    }
    ```
    For private chat you have to add receiver but for group chat there is no need of receiver
5) Retrieve All ChatRoom
    
    ```
    query{
      chatRooms{
        chatRoomID
        creatorID
        creator{
          id
          name
          email
          contact
          bio
          createdAt
        }
        chatRoomName
        chatRoomType
        members{
          id
          chatRoomID
          joinAt
        }
        createdAt
      }
    }
    ```
6) Add Member in Private ChatRoom

    ```
    mutation{
      newChatRoomMember(input:{
        chatRoomID:"848451",
        memberID:"484515"
      },receiverId:"451515"){
        id
        chatRoomID
        joinAt
      }
    }
    ```
7) Add Member in Group ChatRoom

    ```
    mutation{
      newChatRoom(input:{
        creatorID:"15845841451",
        chatRoomName:"abcd",
        chatRoomType:GROUP
      }){
        creatorID
        chatRoomID
        creator{
          name
          email
          contact
          profilePicture
          bio
          createdAt
        }
        chatRoomName
        chatRoomType
        createdAt
      }
    }
    ```
 8) Retrieve ChatConversation by particular charoom
 
    ```
    query{
      chatconversationByChatRoomId(chatRoomID:"0123456789",memberID:"987456310"){
        chatRoomID
        senderId
        sender{
          id
          name
          email
          contact
          profilePicture
          bio
        }
        message
        messageType
        createdAt
        updatedAt
      }
    }
    ```
9) Add Message in ChatRoom 

    ```
    mutation{
      newMessage(input:{
        chatRoomID:"0123456789",
        senderId: "9874563210",
        message:"XXXX",
        messageType:TEXT/IMAGE/VIDEO/GIF,
        messageStatus:SEND
      },senderID:"875456892145"){
        chatRoomID
        senderId
        sender{
          name 
          email
          contact
          profilePicture
        }
        message
        messageType
        messageStatus
        createdAt
      }
    }
    ```
 10) 