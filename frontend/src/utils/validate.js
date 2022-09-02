export function requiredValidate(val, paramInfo, params) {
  if((val != null) && (val !== '')) {
    return true
  }
  if(paramInfo.required === true) {
    if(val == null || val === '') {
      return false
    }
  }else if(paramInfo.required) {
    let paramKeys = Object.keys(paramInfo.required)
    for(let i = 0; i < paramKeys.length; i++) {
      let paramKey = paramKeys[i];
      let paramValuesSet = new Set(paramInfo.required[paramKey])
      let param = params.find(p => p.key === paramKey)
      if(param) {
        if(paramValuesSet.has(param.value)) {
          return false
        }
        if(param.type === 'groupSelect') {
          for(let i = 0; i < param.options.length; i++) {
            let paramOption = param.options[i]
            for(let j = 0; j < paramOption.options.length; j++) {
              if((param.value === paramOption.options[j].value) && paramValuesSet.has(paramOption.category)) {
                return false
              }
            }
          }
        }
        
      }
    }
  }
  return true
}

export function rangeValidate(val, paramInfo) {
  if((paramInfo.min != null) && (val < paramInfo.min)) {
    return false
  }
  if((paramInfo.max != null) && (val > paramInfo.max)) {
    return false
  }
  return true
}

export function requiredMsg(paramInfo) {
  return '该参数必填'
}

export function outRangeMsg(paramInfo) {
  let errorMsg = ''
  if((paramInfo.min != null) && (paramInfo.max != null)) {
    errorMsg = `超过取值范围, 最小值:${paramInfo.min}, 最大值:${paramInfo.max}`
  }else if((paramInfo.min != null) && (paramInfo.max == null)) {
    errorMsg = `超过取值范围, 最小值:${paramInfo.min}`
  }else if((paramInfo.min == null) && (paramInfo.max != null)) {
    errorMsg = `超过取值范围, 最大值:${paramInfo.max}`
  }
  return errorMsg
}