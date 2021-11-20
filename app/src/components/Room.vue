<template>
	<div class="room">
		<h4>
			{{room.name}}
			<i @click="onClick()" :class="[joined ? 'fa-times' : 'fa-door-open', 'fas']"></i>
		</h4>
		<p>Online: {{room.count}}</p>
	</div>
</template>


<script>
	export default{
		name: 'Room',
		data(){
			return{
				joined: false
			}
		},
		props:{
			room: Object	
		},
		components:{
		},
		methods:{
			onClick(){
				this.joined = !this.joined	
			}
		}
	}
</script>

<style scoped>
.fa-door-open{
	color: green;
}
.fa-times{
	color: red;
}
.fas{
	margin: 5px;
}

.room {
  background: #f4f4f4;
  margin: 5px;
  padding: 3px 3px;
}
.room h4 {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.btn {
  display: inline-block;
  background: #000;
  color: #fff;
  border: none;
  padding: 5px 10px;
  margin: 2px;
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

