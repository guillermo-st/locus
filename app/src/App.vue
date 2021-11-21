<template>
	<h1> Locus. </h1>
	<Rooms :rooms="rooms" 
		@refresh-rooms="getRooms()" 
		@add-room="addRoom" 
		@join-room="joinRoom"
		@leave-room="leaveRoom"/>
</template>

<script>
import Rooms from './components/Rooms';
import Button from './components/Button';


export default {
	name: 'App',
  	components: {
  		Button,Rooms
  	},
  	data(){
		return{
			rooms: null,
			wsConn: null
		}
  	},
	methods: {
		getRooms(){
			fetch('http://localhost:8000/rooms')
				.then(response => response.json())
				.then(data => (this.rooms = data));
		},
		addRoom(room){
			fetch("http://localhost:8000/rooms", {
				
				mode: 'no-cors',
			  	method: "POST",
			  	headers: {
			   		'Accept': 'application/json',
			    	'Content-Type': 'application/json'
			  	},
			
			  	body: JSON.stringify(room)
			});
			this.getRooms();
		},
		newWsConn(){
			this.wsConn = new WebSocket("ws://localhost:8000/ws");
	//		this.wsConn.onmessage = this.processMessage(ev);
		},
		processMessage(ev){
			console.log(ev.data);
		},
		joinRoom(roomName){
			const Msg = {
				Action: 'join',
				Room: roomName,
				Username: 'Guille',
			}
			this.wsConn.send(JSON.stringify(Msg));
			this.getRooms();
		},
		leaveRoom(roomName){
			const Msg = {
				Action: 'leave',
				Room: roomName,
				Username: 'Guille',
			}
			this.wsConn.send(JSON.stringify(Msg));
			this.getRooms();
		}
	},
 	created() {
		this.getRooms();
		this.newWsConn();
  	}
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400&display=swap');
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}
body {
  font-family: 'Poppins', sans-serif;
}
.container {
  max-width: 300px;
  margin: 20px auto;
  overflow: auto;
  min-height: 600px;
  border: 1px solid steelblue;
  padding: 30px;
  border-radius: 5px;
}
.btn {
  display: inline-block;
  background: #000;
  color: #fff;
  border: none;
  padding: 10px 20px;
  margin: 5px;
  border-radius: 5px;
  cursor: pointer;
  text-decoration: none;
  font-size: 15px;
  font-family: inherit;
}
.btn:focus {
  outline: none;
}
.btn:active {
  transform: scale(0.98);
}
.btn-block {
  display: block;
  width: 100%;
}
</style>
