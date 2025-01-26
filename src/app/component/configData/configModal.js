import React, { useState } from "react";
import { Button, Modal, Form, Input, Checkbox } from "antd";

export default function ConfigModal() {
    const [open, setOpen] = useState(false);
    const [confirmLoading, setConfirmLoading] = useState(false);

    const showModal = () => {
        setOpen(true);
    };

    const handleOk = () => {
        setConfirmLoading(true);
        setTimeout(() => {
            setOpen(false);
            setConfirmLoading(false);
        }, 2000);
    };

    const handleCancel = () => {
        console.log("Clicked cancel button");
        setOpen(false);
    };

    const onFinish = (values) => {
        console.log("Form Success:", values);
        // 可以在这里处理表单提交逻辑
        handleOk(); // 提交成功后可以关闭模态框
    };

    const onFinishFailed = (errorInfo) => {
        console.log("Form Failed:", errorInfo);
    };

    return (
        <>
            <Button type="primary" onClick={showModal}>
                修改配置
            </Button>
            <Modal
                title="修改配置"
                open={open}
                onOk={handleOk}
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
                                required:true,
                                message:"请输入索引号"
                            }
                        ]}
                    >
                        <Input/>
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
                        <Input/>
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
