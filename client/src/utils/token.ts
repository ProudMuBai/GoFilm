

// 将 token 存储到 localStorage 中
const setToken = (token:string)=>{
    const auth = {key: "auth-token", value: token}
    localStorage.setItem("auth", JSON.stringify(auth))
}

// 获取 localStorage 中的 token 信息
const getToken = ()=>{
    return JSON.parse(localStorage.getItem("auth") as  string)
}

// 删除 localStorage 中的 token 信息
const clearAuthToken = ()=>{
    localStorage.removeItem("auth")
}
export {setToken, getToken,clearAuthToken}