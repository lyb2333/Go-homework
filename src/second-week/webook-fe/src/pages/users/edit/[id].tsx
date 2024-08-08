import React, { useEffect, useState } from 'react';
import { Button, DatePicker, Form, Input } from 'antd';
import axios from "@/axios/axios";
import moment from 'moment';
import router, { useRouter } from "next/router";
import { Profile } from '../model';

const { TextArea } = Input;

const onFinishFailed = (errorInfo: any) => {
    alert("输入有误")
};

function EditForm() {
    let p: Profile = { Email: "", Phone: "", Nickname: "", Birthday: "", AboutMe: "" }
    const [data, setData] = useState<Profile>(p)
    const [isLoading, setLoading] = useState(false)
    const routerInfo = useRouter()
    const { id } = routerInfo.query

    useEffect(() => {
        setLoading(true)
        axios.get('/users/profile', { params: { id } })
            .then((res) => res.data)
            .then((data) => {
                setData(data)
                setLoading(false)
            })
    }, [id])

    const onFinish = (values: any) => {
        values.id = id
        if (values.birthday) {
            values.birthday = moment(values.birthday).format("YYYY-MM-DD")
        }
        axios.post("/users/edit", values)
            .then((res) => {
                if (res.status != 200) {
                    alert(res.statusText);
                    return
                }
                if (res.data?.code == 0) {
                    router.push('/users/profile/' + id)
                    return
                }
                alert(res.data?.msg || "系统错误");
            }).catch((err) => {
                alert(err);
            })
    };

    if (isLoading) return <p>Loading...</p>
    if (!data) return <p>No profile data</p>
    return <Form
        name="basic"
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 16 }}
        style={{ maxWidth: 600 }}
        initialValues={{
            birthday: moment(data.Birthday, 'YYYY-MM-DD'),
            nickname: data.Nickname,
            aboutMe: data.AboutMe
        }}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        autoComplete="off"
    >
        <Form.Item
            label="昵称"
            name="nickname"
        >
            <Input />
        </Form.Item>

        <Form.Item
            label="生日"
            name="birthday"
        >
            <DatePicker format={"YYYY-MM-DD"}
                placeholder={""} />
        </Form.Item>

        <Form.Item
            label="关于我"
            name="aboutMe"
        >
            <TextArea rows={4} />
        </Form.Item>

        <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
            <Button type="primary" htmlType="submit">
                提交
            </Button>
        </Form.Item>
    </Form>
}

export default EditForm;