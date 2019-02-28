# Schema Types

<details>
  <summary><strong>Table of Contents</strong></summary>

  * [Query](#query)
  * [Mutation](#mutation)
  * [Objects](#objects)
    * [ChatConversation](#chatconversation)
    * [ChatRoom](#chatroom)
    * [ChatRoomList](#chatroomlist)
    * [Member](#member)
    * [MemberCountsWithMemberDetailsByChatRoom](#membercountswithmemberdetailsbychatroom)
    * [Subscription](#subscription)
    * [User](#user)
  * [Inputs](#inputs)
    * [DeleteChat](#deletechat)
    * [DeleteMessage](#deletemessage)
    * [LeaveChatRoom](#leavechatroom)
    * [NewChatRoomMembers](#newchatroommembers)
    * [NewGroupChatRoom](#newgroupchatroom)
    * [NewMessage](#newmessage)
    * [NewPrivateChatRoom](#newprivatechatroom)
    * [NewUser](#newuser)
    * [RemoveChatRoomPermenantly](#removechatroompermenantly)
    * [RemoveMembersFromChatRoom](#removemembersfromchatroom)
    * [UpdateChatRoomDetail](#updatechatroomdetail)
    * [UpdateMessage](#updatemessage)
    * [UpdateMessageStatus](#updatemessagestatus)
    * [UpdateUser](#updateuser)
  * [Enums](#enums)
    * [ChatRoomType](#chatroomtype)
    * [MessageType](#messagetype)
    * [State](#state)
  * [Scalars](#scalars)
    * [Boolean](#boolean)
    * [Float](#float)
    * [ID](#id)
    * [Int](#int)
    * [String](#string)
    * [Time](#time)

</details>

## Query
<table>
<thead>
<tr>
<th align="left">Field</th>
<th align="right">Argument</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>users</strong></td>
<td valign="top">[<a href="#user">User</a>!]!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">name</td>
<td valign="top"><a href="#string">String</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRooms</strong></td>
<td valign="top">[<a href="#chatroom">ChatRoom</a>!]</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatconversationByChatRoomId</strong></td>
<td valign="top">[<a href="#chatconversation">ChatConversation</a>!]!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">chatRoomID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">memberID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>memberListByChatRoomId</strong></td>
<td valign="top"><a href="#membercountswithmemberdetailsbychatroom">MemberCountsWithMemberDetailsByChatRoom</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">chatRoomID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">memberID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomListByMemberId</strong></td>
<td valign="top">[<a href="#chatroomlist">ChatRoomList</a>!]!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">memberID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>MemberLogIn</strong></td>
<td valign="top"><a href="#user">User</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">name</td>
<td valign="top"><a href="#string">String</a>!</td>
<td></td>
</tr>
</tbody>
</table>

## Mutation
<table>
<thead>
<tr>
<th align="left">Field</th>
<th align="right">Argument</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>newUser</strong></td>
<td valign="top"><a href="#user">User</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#newuser">NewUser</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>newPrivateChatRoom</strong></td>
<td valign="top"><a href="#chatroom">ChatRoom</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#newprivatechatroom">NewPrivateChatRoom</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>newGroupchatRoom</strong></td>
<td valign="top"><a href="#chatroom">ChatRoom</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#newgroupchatroom">NewGroupChatRoom</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>newChatRoomMembers</strong></td>
<td valign="top"><a href="#boolean">Boolean</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#newchatroommembers">NewChatRoomMembers</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>newMessage</strong></td>
<td valign="top"><a href="#chatconversation">ChatConversation</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#newmessage">NewMessage</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>updateMessage</strong></td>
<td valign="top"><a href="#chatconversation">ChatConversation</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#updatemessage">UpdateMessage</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>updateMessageStatus</strong></td>
<td valign="top"><a href="#chatconversation">ChatConversation</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#updatemessagestatus">UpdateMessageStatus</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>deleteMessage</strong></td>
<td valign="top"><a href="#chatconversation">ChatConversation</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#deletemessage">DeleteMessage</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>deleteChatRoomPermenantly</strong></td>
<td valign="top"><a href="#chatroom">ChatRoom</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#removechatroompermenantly">RemoveChatRoomPermenantly</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>removeMembersFromChatRoomByCreator</strong></td>
<td valign="top"><a href="#chatroom">ChatRoom</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#removemembersfromchatroom">RemoveMembersFromChatRoom</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>leaveChatRoom</strong></td>
<td valign="top"><a href="#boolean">Boolean</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#leavechatroom">LeaveChatRoom</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>updateChatRoomDetail</strong></td>
<td valign="top"><a href="#chatroom">ChatRoom</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#updatechatroomdetail">UpdateChatRoomDetail</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>deleteChatRoom</strong></td>
<td valign="top"><a href="#boolean">Boolean</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">input</td>
<td valign="top"><a href="#deletechat">DeleteChat</a>!</td>
<td></td>
</tr>
</tbody>
</table>

## Objects

### ChatConversation

<table>
<thead>
<tr>
<th align="left">Field</th>
<th align="right">Argument</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>messageId</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>senderID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>sender</strong></td>
<td valign="top"><a href="#user">User</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>message</strong></td>
<td valign="top"><a href="#string">String</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageType</strong></td>
<td valign="top"><a href="#messagetype">MessageType</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageStatus</strong></td>
<td valign="top"><a href="#state">State</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageParentId</strong></td>
<td valign="top"><a href="#id">ID</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>createdAt</strong></td>
<td valign="top"><a href="#time">Time</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>updatedAt</strong></td>
<td valign="top"><a href="#time">Time</a></td>
<td></td>
</tr>
</tbody>
</table>

### ChatRoom

<table>
<thead>
<tr>
<th align="left">Field</th>
<th align="right">Argument</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>creatorID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>creator</strong></td>
<td valign="top"><a href="#user">User</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomName</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomType</strong></td>
<td valign="top"><a href="#chatroomtype">ChatRoomType</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>members</strong></td>
<td valign="top">[<a href="#member">Member</a>!]!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>createdAt</strong></td>
<td valign="top"><a href="#time">Time</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>updateByID</strong></td>
<td valign="top"><a href="#id">ID</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>updateBy</strong></td>
<td valign="top"><a href="#user">User</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>updatedAt</strong></td>
<td valign="top"><a href="#time">Time</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>deleteAt</strong></td>
<td valign="top"><a href="#time">Time</a></td>
<td></td>
</tr>
</tbody>
</table>

### ChatRoomList

<table>
<thead>
<tr>
<th align="left">Field</th>
<th align="right">Argument</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>name</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomType</strong></td>
<td valign="top"><a href="#string">String</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>createdAt</strong></td>
<td valign="top"><a href="#time">Time</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>totalMember</strong></td>
<td valign="top"><a href="#int">Int</a></td>
<td></td>
</tr>
</tbody>
</table>

### Member

<table>
<thead>
<tr>
<th align="left">Field</th>
<th align="right">Argument</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>id</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>member</strong></td>
<td valign="top"><a href="#user">User</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>joinAt</strong></td>
<td valign="top"><a href="#time">Time</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>deleteAt</strong></td>
<td valign="top"><a href="#time">Time</a></td>
<td></td>
</tr>
</tbody>
</table>

### MemberCountsWithMemberDetailsByChatRoom

<table>
<thead>
<tr>
<th align="left">Field</th>
<th align="right">Argument</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>memberCount</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>members</strong></td>
<td valign="top">[<a href="#member">Member</a>!]!</td>
<td></td>
</tr>
</tbody>
</table>

### Subscription

<table>
<thead>
<tr>
<th align="left">Field</th>
<th align="right">Argument</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>userJoined</strong></td>
<td valign="top"><a href="#user">User</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messagePost</strong></td>
<td valign="top"><a href="#chatconversation">ChatConversation</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">chatRoomID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageStatusUpdate</strong></td>
<td valign="top"><a href="#chatconversation">ChatConversation</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">messageID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">chatRoomID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageUpdate</strong></td>
<td valign="top"><a href="#chatconversation">ChatConversation</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">chatRoomID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageDelete</strong></td>
<td valign="top"><a href="#chatconversation">ChatConversation</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">chatRoomID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomListByMember</strong></td>
<td valign="top">[<a href="#chatroomlist">ChatRoomList</a>!]!</td>
<td></td>
</tr>
<tr>
<td colspan="2" align="right" valign="top">memberID</td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
</tbody>
</table>

### User

<table>
<thead>
<tr>
<th align="left">Field</th>
<th align="right">Argument</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>id</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>userName</strong></td>
<td valign="top"><a href="#string">String</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>firstName</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>lastName</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>email</strong></td>
<td valign="top"><a href="#string">String</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>contact</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>bio</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>profilePicture</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>createdAt</strong></td>
<td valign="top"><a href="#time">Time</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>updatedAt</strong></td>
<td valign="top"><a href="#time">Time</a></td>
<td></td>
</tr>
</tbody>
</table>

## Inputs

### DeleteChat

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>memberID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
</tbody>
</table>

### DeleteMessage

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>DeleteByID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
</tbody>
</table>

### LeaveChatRoom

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>memberID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
</tbody>
</table>

### NewChatRoomMembers

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>memberIDs</strong></td>
<td valign="top">[<a href="#id">ID</a>!]</td>
<td></td>
</tr>
</tbody>
</table>

### NewGroupChatRoom

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>creatorID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomName</strong></td>
<td valign="top"><a href="#string">String</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomType</strong></td>
<td valign="top"><a href="#chatroomtype">ChatRoomType</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>receiverID</strong></td>
<td valign="top">[<a href="#id">ID</a>!]</td>
<td></td>
</tr>
</tbody>
</table>

### NewMessage

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>senderID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>message</strong></td>
<td valign="top"><a href="#string">String</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageType</strong></td>
<td valign="top"><a href="#messagetype">MessageType</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageParentId</strong></td>
<td valign="top"><a href="#id">ID</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageStatus</strong></td>
<td valign="top"><a href="#state">State</a>!</td>
<td></td>
</tr>
</tbody>
</table>

### NewPrivateChatRoom

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>creatorID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomType</strong></td>
<td valign="top"><a href="#chatroomtype">ChatRoomType</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>receiverID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
</tbody>
</table>

### NewUser

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>userName</strong></td>
<td valign="top"><a href="#string">String</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>firstName</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>lastName</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>email</strong></td>
<td valign="top"><a href="#string">String</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>contact</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>profilePicture</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>bio</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
</tbody>
</table>

### RemoveChatRoomPermenantly

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>creator</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
</tbody>
</table>

### RemoveMembersFromChatRoom

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>creatorID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>memberID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
</tbody>
</table>

### UpdateChatRoomDetail

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomName</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>updateByID</strong></td>
<td valign="top"><a href="#id">ID</a></td>
<td></td>
</tr>
</tbody>
</table>

### UpdateMessage

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>message</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>senderID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>messageID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>chatRoomID</strong></td>
<td valign="top"><a href="#id">ID</a>!</td>
<td></td>
</tr>
</tbody>
</table>

### UpdateMessageStatus

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>messageStatus</strong></td>
<td valign="top"><a href="#state">State</a>!</td>
<td></td>
</tr>
</tbody>
</table>

### UpdateUser

<table>
<thead>
<tr>
<th colspan="2" align="left">Field</th>
<th align="left">Type</th>
<th align="left">Description</th>
</tr>
</thead>
<tbody>
<tr>
<td colspan="2" valign="top"><strong>userName</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>firstName</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
<tr>
<td colspan="2" valign="top"><strong>lastName</strong></td>
<td valign="top"><a href="#string">String</a></td>
<td></td>
</tr>
</tbody>
</table>

## Enums

### ChatRoomType

<table>
<thead>
<th align="left">Value</th>
<th align="left">Description</th>
</thead>
<tbody>
<tr>
<td valign="top"><strong>PRIVATE</strong></td>
<td></td>
</tr>
<tr>
<td valign="top"><strong>GROUP</strong></td>
<td></td>
</tr>
</tbody>
</table>

### MessageType

<table>
<thead>
<th align="left">Value</th>
<th align="left">Description</th>
</thead>
<tbody>
<tr>
<td valign="top"><strong>TEXT</strong></td>
<td></td>
</tr>
<tr>
<td valign="top"><strong>IMAGE</strong></td>
<td></td>
</tr>
<tr>
<td valign="top"><strong>VIDEO</strong></td>
<td></td>
</tr>
<tr>
<td valign="top"><strong>GIF</strong></td>
<td></td>
</tr>
<tr>
<td valign="top"><strong>AUDIO</strong></td>
<td></td>
</tr>
</tbody>
</table>

### State

<table>
<thead>
<th align="left">Value</th>
<th align="left">Description</th>
</thead>
<tbody>
<tr>
<td valign="top"><strong>SEND</strong></td>
<td></td>
</tr>
<tr>
<td valign="top"><strong>UNREAD</strong></td>
<td></td>
</tr>
<tr>
<td valign="top"><strong>READ</strong></td>
<td></td>
</tr>
<tr>
<td valign="top"><strong>DELETE</strong></td>
<td></td>
</tr>
<tr>
<td valign="top"><strong>ADD</strong></td>
<td></td>
</tr>
</tbody>
</table>

## Scalars

### Boolean

### Float

### ID

### Int

### String

### Time

