"use client";
import React, { useState } from "react";
import {Button, Checkbox, Form, Input, message, Space, Typography} from "antd";

const { Title } = Typography;

export default function MyLogin() {
    const [loading, setLoading] = useState(false); // 用于按钮加载状态

    const [messageApi, contextHolder] = message.useMessage();

    const onFinish = async (values) => {
        console.log("Success:", values);
        setLoading(true);
        try {
            const response = await fetch("http://localhost:8080/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(values),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();

            if (result.success) {
                localStorage.setItem("token", result.token);
                // 登录成功，跳转到目标页面
                // message.success("Login successful!");
                // messageApi.info("Login successful!");
                window.location.href = "http://192.168.56.1:3000/data-show";
            } else {
                // 登录失败，显示错误消息
                messageApi.info(result.message || "Username or password is incorrect");
                // message.error(result.message || "Username or password is incorrect");
            }
        } catch (error) {
            console.error("Error:", error);
            // message.error("Login failed. Please try again later.");
            messageApi.info("Login failed. Please try again later.");
        } finally {
            setLoading(false);
        }
    };

    const onFinishFailed = (errorInfo) => {
        console.log("Failed:", errorInfo);
        message.error("Please complete the form before submitting.");
    };

    return (
        <Space direction="vertical" size="middle"
               style={{
            display: "flex", // 使用 flex 布局
            justifyContent: "center", // 水平居中
            alignItems: "center", // 垂直居中
            height: "100vh", // 父容器高度占满整个视口
            background: "white", // 背景色，可选
               }}>
        {contextHolder}
            <div style={{
                background: "rgba(0, 0, 0, 0.03)",
                borderRadius: 8,}}>
        <Title level={2} type="success" style={{
            textAlign:"center"
        }}>IDs后端管理系统</Title>
        <Form
            name="basic"
            labelCol={{
                span: 5,
            }}
            wrapperCol={{
                span: 16,
            }}
            style={{
                maxWidth: 600,
                minWidth: 400,
            }}
            initialValues={{
                remember: true,
            }}
            onFinish={onFinish}
            onFinishFailed={onFinishFailed}
            autoComplete="off"
        >
            <Form.Item
                label="Username"
                name="username"
                rules={[
                    {
                        required: true,
                        message: "Please input your username!",
                    },
                ]}
            >
                <Input />
            </Form.Item>

            <Form.Item
                label="Password"
                name="password"
                rules={[
                    {
                        required: true,
                        message: "Please input your password!",
                    },
                ]}
            >
                <Input.Password />
            </Form.Item>

            <Form.Item name="remember" valuePropName="checked" label={null}>
                <Checkbox>Remember me</Checkbox>
            </Form.Item>

            <Form.Item label={null}>
                <Button type="primary" htmlType="submit" loading={loading}>
                    Submit
                </Button>
            </Form.Item>
        </Form>
            </div>
        </Space>
    );
}
