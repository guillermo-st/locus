<template>
	<div class="room">
		<h3>
			{{room.name}}
			<Button label='Join' color="green"/>
		</h3>
		<p>Online: {{room.count}}</p>
	</div>
</template>


<script>
	import Button from './Button';
	export default{
		name: 'Room',
		props:{
			room: Object	
		},
		components:{
			Button
		},
		methods:{
			onClick(){
				console.log('click')
			}
		}
	}
</script>

<style scoped>

.room {
  background: #f4f4f4;
  margin: 5px;
  padding: 3px 3px;
}
.room h3 {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

</style>

