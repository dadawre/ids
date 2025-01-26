"use client";
import React, { useState } from "react";
import { Button, Modal, Form, Input, message } from "antd";

export default function ConfigModal({data,onChange}) {
    const [open, setOpen] = useState(false);
    const [confirmLoading, setConfirmLoading] = useState(false);

    const handleChangeConfigData = (changeData)=>{
        const nextData = data.map((v,i)=>{
            if (v.index === changeData.index){
                changeData.key = v.key
                return changeData
            }else{
                return v
            }
        })
        onChange(nextData)
    }

    const showModal = () => {
        setOpen(true);
    };

    const handleOk = () => {
        setOpen(false);
    };

    const handleCancel = () => {
        console.log("Clicked cancel button");
        setOpen(false);
    };

    const onFinish = async (values) => {

        const requestData = {
            index: parseInt(values.index, 10), // 将 index 转为整数
            value: values.value
        };

        console.log("Form Success:", requestData);
        console.log("Sending data to server:", JSON.stringify(requestData)); // 打印发送的数据
        setConfirmLoading(true);
        try {
            const response = await fetch("http://localhost:8080/changeConfig", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(requestData),
            });

            if (!response.ok) {
                message.error("请求错误，请重试");
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();

            if (result.success) {
                message.success("配置修改成功");
                handleChangeConfigData(requestData)
                handleOk(); // 关闭对话框
            } else {
                message.error(result.message || "配置修改失败");
            }
        } catch (error) {
            console.error("Error:", error);
            message.error("配置修改失败，请稍后再试");
        } finally {
            setConfirmLoading(false);
        }
    };

    const onFinishFailed = (errorInfo) => {
        console.log("Form Failed:", errorInfo);
        message.error("请填写完整表单后再提交");
    };

    return (
        <>
            <Button type="primary" onClick={showModal}>
                修改配置
            </Button>
            <Modal
                title="修改配置"
                open={open}
                confirmLoading={confirmLoading}
                onCancel={handleCancel}
                footer={null} // 去掉默认的 Footer，用表单按钮替代
            >
                <Form
                    name="basic"
                    labelCol={{
                        span: 8,
                    }}
                    wrapperCol={{
                        span: 16,
                    }}
                    onFinish={onFinish}
                    onFinishFailed={onFinishFailed}
                    autoComplete="off"
                >
                    <Form.Item
                        label="索引"
                        name="index"
                        rules={[
                            {
                                required: true,
                                message: "请输入索引号",
                            },
                        ]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item
                        label="配置值"
                        name="value"
                        rules={[
                            {
                                required: true,
                                message: "请输入配置值",
                            },
                        ]}
                    >
                        <Input />
                    </Form.Item>

                    <Form.Item
                        wrapperCol={{
                            offset: 8,
                            span: 16,
                        }}
                    >
                        <Button type="primary" htmlType="submit">
                            提交
                        </Button>
                    </Form.Item>
                </Form>
            </Modal>
        </>
    );
}
