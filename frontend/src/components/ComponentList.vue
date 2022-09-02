<template>
  <div class="component-list">
    <div class="component-list-header">
      <a-input-search 
        v-model:value="searchVal"
        placeholder="搜索组件..."
        style="width: 100%"
        allowClear
        @change="searchDebounce"
      >
      </a-input-search>
    </div>
    <div class="component-list-content">
      <MLTree :treeData="componentData" @start-drag="startDragComponent"></MLTree>
    </div>
  </div>
</template>

<script>
import MLTree from './tree/MLTree.vue'
import { deepCopy, debounce }  from '../utils/index'

export default {
  name: 'ComponentList',
  components: {
    MLTree
  },
  emits: ['start-drag', 'dragging', 'end-drag'],
  data() {
    return {
      componentData: [],
      searchVal: ''
    }
  },
  created() {
    this.searchDebounce = debounce(this.search, 100)
  },
  watch: {
    '$store.state.graph.componentRawData': {
      handler() {
        this.allData = this.formatData(deepCopy(this.$store.state.graph.componentRawData))
        this.search()
      }
    }
  },
  methods: {
    formatData(rawData) {
      if(!rawData) {
        return []
      }

      // 过滤掉不用显示的组件
      rawData = rawData.filter(comp => {
        return (comp.key !== 'ProcessIn') &&  (comp.key !== 'ProcessOut')
      })

      let typeMap = {};
      for(let i = 0, len = rawData.length; i < len; i++) {
        let rawComp = rawData[i]
        rawComp.title = rawComp.name;
        if(typeMap[rawComp.type] == null) {
          typeMap[rawComp.type] = {
            title: rawComp.typeLabel,
            key: rawComp.type,
            children: []
          }
          if(rawComp.category == null) {
            typeMap[rawComp.type].children.push(rawComp)
          }else {
            typeMap[rawComp.type].children.push({
              title: rawComp.categoryLabel,
              key: rawComp.category,
              icon: rawComp.icon,
              children: [rawComp]
            })
          }
        }else {
          if(rawComp.category == null) {
            typeMap[rawComp.type].children.push(rawComp)
          }else {
            let category = typeMap[rawComp.type].children.find(c => c.key === rawComp.category)
            if(category) {
              category.children.push(rawComp)
            }else {
              typeMap[rawComp.type].children.push({
                title: rawComp.categoryLabel,
                key: rawComp.category,
                children: [rawComp],
              })
            }
          }
        }
      }

      return Object.values(typeMap);
    },
    startDragComponent(componentData, x, y) {
      this.$emit('start-drag', componentData, x, y)
      window.addEventListener('mouseup', this.dropComponent)
      window.addEventListener('mousemove', this.dragginComponent)
    },
    dragginComponent(e) {
      this.$emit('dragging', e.clientX, e.clientY)
    },
    dropComponent(e) {
      this.$emit('end-drag', e)
      window.removeEventListener('mousemove', this.dragginComponent)
      window.removeEventListener('mouseup', this.dropComponent)
    },
    search() {
      let val = this.searchVal.trim()
      if(!val) {
        this.componentData = this.allData
      }else {
        let resData = []
        val = val.toUpperCase()
        for(let i = 0, len = this.allData.length; i < len; i++) {
          let children = this.allData[i].children
          this.allData[i].children = []
          let newItem = deepCopy(this.allData[i])
          this.allData[i].children = children

          if(this.isSearched(val, this.allData[i], newItem)) {
            resData.push(newItem)
          }
        }
        this.expendAll(resData)
        this.componentData = resData
      }
    },
    isSearched(searchVal, sourceItem, targetItem) {
      if(targetItem.isDirectory) {
        if(targetItem.title.toUpperCase().indexOf(searchVal) !== -1) {
          targetItem.children = deepCopy(sourceItem.children)
          return true
        }else {
          for(let i = 0, len = sourceItem.children.length; i < len; i++) {
            let children = sourceItem.children[i].children
            sourceItem.children[i].children = []
            let newItem = deepCopy(sourceItem.children[i])
            sourceItem.children[i].children = children
            if(this.isSearched(searchVal, sourceItem.children[i], newItem)) {
              targetItem.children.push(newItem)
            }
          }
          return targetItem.children.length > 0
        }
      }else {
        return targetItem.title.toUpperCase().indexOf(searchVal) !== -1
      }
    },
    expendAll(treeDatas) {
      for(let i = 0, len = treeDatas.length; i < len; i++) {
        treeDatas[i].expand = true
        if(treeDatas[i].children) {
          this.expendAll(treeDatas[i].children)
        }
      }
    }
  }
}
</script>

<style lang="less">

</style>