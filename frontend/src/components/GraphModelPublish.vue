<template>
  <a-modal 
    v-model:visible="modalVisible" 
    title="模型发布"
    :maskClosable="false"
    :footer="null"
    :afterClose="afterCloseHandler">
      <a-form
        ref="model"
        :model="modelState"
        :rules="rules"
        :label-col="labelCol"
        :wrapper-col="wrapperCol"
      >
      <a-form-item label="发布类型" name="type">
        <a-radio-group v-model:value="modelState.type">
          <a-radio :value="0">新建</a-radio>
          <a-radio :value="1">覆盖</a-radio>
        </a-radio-group>
      </a-form-item>
      <a-form-item v-if="modelState.type === 1" label="覆盖模型" name="id">
        <a-select
          v-model:value="modelState.id"
          show-search
          :filter-option="false"
          @search="handleSearch"
          placeholder="请选择一个要覆盖的模型"
          style="width: 100%"
          :loading="modelListLoading"
          :options="searchedModelList"
        >
        </a-select>
      </a-form-item>
      <a-form-item label="名称" name="name">
        <a-input v-model:value="modelState.name" />
      </a-form-item>
      <a-form-item label="描述" name="description">
        <a-input v-model:value="modelState.description" />
      </a-form-item>
      <a-form-item :wrapper-col="{ span: 4, offset: 19 }">
        <a-button type="primary" @click="onSubmit">创建</a-button>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script>
import { getModelList } from '../service'
import { debounce } from '../utils' 

export default {
  name: 'GraphModelPublish',
  props: {
    visible: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:visible', 'confirm'],
  data() {
    return {
      modalVisible: this.visible,
      modelListLoading: false,
      modelList: [],
      searchVal: '',
      modelState: {
        type: 0,   // 0:新建 1：覆盖
        id: null,
        name: '',
        description: ''
      },
      rules: {
        name: [
           { required: true, message: '请输入模型名称', trigger: 'blur' }
        ],
        id: [
          { required: true, message: '请选择一个要覆盖的模型', trigger: 'change' }
        ]
      },
      labelCol: { span: 4 },
      wrapperCol: { span: 18 },
    }
  },
  watch: {
    visible() {
      this.modalVisible = this.visible
    },
    'modelState.type': {
      handler() {
        if((this.modelState.type === 1) && (this.modelList.length === 0)) {
          this.getModels()
        }
      }
    }
  },
  created() {},
  computed: {
    searchedModelList() {
      if(this.searchVal) {
        return this.modelList.filter(model => model.label.toLowerCase().indexOf(this.searchVal.toLowerCase()) > -1);
      }else {
        return this.modelList
      }
    }
  },
  methods: {
    getModels() {
      this.modelListLoading = true
      getModelList().then(data => {
        this.modelList = data.map(model => {
          return {
            value: model.id,
            label: model.name
          }
        })
      }).catch(err => {
        console.error('get model list error:', err)
      }).finally(() => {
        this.modelListLoading = false
      })
    },
    afterCloseHandler() {
      this.$emit('update:visible', false)
    },
    onSubmit() {
      this.$refs.model.validate().then(() => {
        this.modalVisible = false
        let model = {
          name: this.modelState.name,
          description: this.modelState.description
        };
        if(this.modelState.type === 1) {
          model.id = this.modelState.id
        }
        this.$emit('confirm', model)
      }).catch(() => {})
    },
    handleSearch: debounce(function(val) {
      this.searchVal = val
    }, 200)
  }
}
</script>

<style>

</style>