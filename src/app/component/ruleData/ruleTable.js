"use client";
import "./EditableTable.css";
import React, { useEffect, useState } from "react";
import {
    Button,
    Form,
    Input,
    InputNumber,
    Modal,
    Popconfirm,
    Table,
    Typography,
    message,
} from "antd";
import axios from "axios";

const EditableCell = ({
                          editing,
                          dataIndex,
                          title,
                          inputType,
                          record,
                          index,
                          children,
                          ...restProps
                      }) => {
    let inputNode;

    if (dataIndex === "StringParameterStrings") {
        inputNode = <Input.TextArea rows={2} />;
    } else {
        inputNode = inputType === "number" ? <InputNumber /> : <Input />;
    }

    return (
        <td {...restProps}>
            {editing ? (
                <Form.Item
                    name={dataIndex}
                    style={{
                        margin: 0,
                    }}
                    rules={[
                        {
                            required: true,
                            message: `Please Input ${title}!`,
                        },
                    ]}
                >
                    {inputNode}
                </Form.Item>
            ) : (
                children
            )}
        </td>
    );
};

export default function RuleTable() {
    const [form] = Form.useForm();
    const [dataSource, setDataSource] = useState([]);
    const [editingKey, setEditingKey] = useState("");
    const [modalOpen, setModalOpen] = useState(false);
    const [modalLoading, setModalLoading] = useState(false);
    const [addForm] = Form.useForm();

    const isEditing = (record) => record.ID.toString() === editingKey;

    useEffect(() => {
        // Fetch initial data from the server
        axios.get("http://localhost:8080/RuleFetch").then((response) => {
            const data = response.data.results.map((item) => ({
                ...item,
                key: item.ID, // 用于表格的唯一标识
            }));
            setDataSource(data);
        });
    }, []);

    const edit = (record) => {
        form.setFieldsValue({ ...record });
        setEditingKey(record.ID.toString());
    };

    const cancel = () => {
        setEditingKey("");
    };

    const save = async (key) => {
        try {
            const row = await form.validateFields();

            // 如果是 StringParameterStrings 字段，将逗号分隔字符串解析为数组
            if (row.StringParameterStrings && typeof row.StringParameterStrings === "string") {
                row.StringParameterStrings = row.StringParameterStrings.split(",").map((item) => item.trim());
            }

            const updatedRow = { ...dataSource.find((item) => item.ID === key), ...row };

            // 调用编辑接口
            await axios.post("http://localhost:8080/RuleEdit", updatedRow);

            const newData = [...dataSource];
            const index = newData.findIndex((item) => item.ID === key);
            if (index > -1) {
                newData.splice(index, 1, updatedRow);
                setDataSource(newData);
                setEditingKey("");
                message.success("Row updated successfully");
            } else {
                message.error("Failed to find the row to update");
            }
        } catch (errInfo) {
            console.error("Save failed:", errInfo);
            message.error("Failed to save changes. Please try again.");
        }
    };

    const handleDelete = (key) => {
        axios.post("http://localhost:8080/RuleDelete", { ID: key }).then(() => {
            const newData = dataSource.filter((item) => item.ID !== key);
            setDataSource(newData);
            message.success("Row deleted successfully");
        });
    };

    const showAddModal = () => {
        setModalOpen(true);
    };

    const handleAdd = () => {
        addForm
            .validateFields()
            .then((values) => {
                setModalLoading(true);

                const newData = {
                    ID: dataSource.length + 1,
                    ProtocolsString: values.ProtocolsString,
                    Priority: values.Priority,
                    ActionString: values.ActionString,
                    MsgString: values.MsgString,
                    StringParameterStrings: values.StringParameterStrings
                        ? values.StringParameterStrings.split(",").map((item) => item.trim())
                        : [],
                };

                axios
                    .post("http://localhost:8080/RuleAdd", newData)
                    .then(() => {
                        setDataSource([...dataSource, newData]);
                        message.success("Row added successfully");
                        setModalOpen(false);
                        setModalLoading(false);
                        addForm.resetFields();
                    })
                    .catch(() => {
                        message.error("Failed to add row. Please try again.");
                        setModalLoading(false);
                    });
            })
            .catch((error) => {
                console.error("Validation failed:", error);
            });
    };

    const handleCancel = () => {
        setModalOpen(false);
        addForm.resetFields();
    };

    const columns = [
        {
            title: "ID",
            dataIndex: "ID",
            width: "5%",
        },
        {
            title: "ProtocolsString",
            dataIndex: "ProtocolsString",
            editable: true,
        },
        {
            title: "Priority",
            dataIndex: "Priority",
            editable: true,
        },
        {
            title: "ActionString",
            dataIndex: "ActionString",
            editable: true,
        },
        {
            title: "MsgString",
            dataIndex: "MsgString",
            editable: true,
        },
        {
            title: "StringParameterStrings",
            dataIndex: "StringParameterStrings",
            render: (text) => (text && text.length > 0 ? text.join(", ") : "N/A"),
            editable: true,
        },
        {
            title: "Operation",
            dataIndex: "operation",
            render: (_, record) => {
                const editable = isEditing(record);
                return editable ? (
                    <span>
            <Typography.Link onClick={() => save(record.ID)} style={{ marginRight: 8 }}>
              Save
            </Typography.Link>
            <Popconfirm title="Sure to cancel?" onConfirm={cancel}>
              <a>Cancel</a>
            </Popconfirm>
          </span>
                ) : (
                    <span>
            <Typography.Link disabled={editingKey !== ""} onClick={() => edit(record)}>
              Edit
            </Typography.Link>
            <Popconfirm title="Sure to delete?" onConfirm={() => handleDelete(record.ID)}>
              <a style={{ marginLeft: 16 }}>Delete</a>
            </Popconfirm>
          </span>
                );
            },
        },
    ];

    const mergedColumns = columns.map((col) => {
        if (!col.editable) {
            return col;
        }
        return {
            ...col,
            onCell: (record) => ({
                record,
                editable: col.editable,
                dataIndex: col.dataIndex,
                title: col.title,
                inputType: typeof record[col.dataIndex] === "number" ? "number" : "text",
                editing: isEditing(record),
            }),
        };
    });

    return (
        <div>
            <Button onClick={showAddModal} type="primary" style={{ marginBottom: 16 }}>
                Add a row
            </Button>
            <Modal
                title="Add New Row"
                open={modalOpen}
                confirmLoading={modalLoading}
                onOk={handleAdd}
                onCancel={handleCancel}
            >
                <Form form={addForm} layout="vertical">
                    <Form.Item name="ProtocolsString" label="ProtocolsString" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="Priority" label="Priority" rules={[{ required: true }]}>
                        <InputNumber style={{ width: "100%" }} />
                    </Form.Item>
                    <Form.Item name="ActionString" label="ActionString" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="MsgString" label="MsgString" rules={[{ required: true }]}>
                        <Input />
                    </Form.Item>
                    <Form.Item name="StringParameterStrings" label="StringParameterStrings">
                        <Input.TextArea rows={2} />
                    </Form.Item>
                </Form>
            </Modal>
            <Form form={form} component={false}>
                <Table
                    components={{
                        body: {
                            cell: EditableCell,
                        },
                    }}
                    bordered
                    dataSource={dataSource}
                    columns={mergedColumns}
                    rowClassName="editable-row"
                    pagination={{
                        onChange: cancel,
                    }}
                    rowKey="ID"
                />
            </Form>
        </div>
    );
}
