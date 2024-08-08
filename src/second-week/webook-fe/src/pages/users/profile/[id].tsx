import { ProDescriptions } from '@ant-design/pro-components';
import React, { useState, useEffect } from 'react';
import { Button } from 'antd';
import axios from "@/axios/axios";
import { useRouter } from 'next/router';
import { Profile } from '../model';

function Page() {
    let p: Profile = { Email: "", Phone: "", Nickname: "", Birthday: "", AboutMe: "" }
    const [data, setData] = useState<Profile>(p)
    const [isLoading, setLoading] = useState(false)
    const router = useRouter()
    const { id } = router.query

    useEffect(() => {
        setLoading(true)
        axios.get('/users/profile', {
            params: {
                id
            }
        })
            .then((res) => res.data)
            .then((data) => {
                setData(data)
                setLoading(false)
            })
    }, [id])

    if (isLoading) return <p>Loading...</p>
    if (!data) return <p>No profile data</p>

    return (
        <ProDescriptions
            column={1}
            title="个人信息"
        >
            <ProDescriptions.Item label="昵称" valueType="text">
                {data.Nickname}
            </ProDescriptions.Item>
            <ProDescriptions.Item
                // span={1}
                valueType="text"
                label="邮箱"
            >{data.Email}
            </ProDescriptions.Item>
            <ProDescriptions.Item
                // span={1}
                valueType="text"
                label="手机"
            >{data.Phone}
            </ProDescriptions.Item>
            <ProDescriptions.Item label="生日" valueType="date">
                {data.Birthday}
            </ProDescriptions.Item>
            <ProDescriptions.Item
                valueType="text"
                label="关于我"
            >
                {data.Aboutme}
            </ProDescriptions.Item>
            <ProDescriptions.Item>
                <Button href={"/users/edit/" + id} type={"primary"}>修改</Button>
            </ProDescriptions.Item>

        </ProDescriptions>
    )
}

export default Page