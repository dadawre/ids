"use client";
import React, { useEffect, useState, useRef } from 'react';
import { Table, Button, Form, Input, Switch, Select, message } from 'antd';
import qs from 'qs';
import { SearchOutlined } from '@ant-design/icons';
import Highlighter from 'react-highlight-words';

const { Option } = Select;

export default function WBTable() {
    // 用于存储表格数据
    const [data, setData] = useState([]);
    // 加载状态
    const [loading, setLoading] = useState(false);
    // 表格参数，包含分页、过滤、排序等
    const [tableParams, setTableParams] = useState({
        pagination: {
            current: 1,
            pageSize: 10,
        },
        filters: {},
        sortOrder: undefined,
        sortField: undefined,
    });
    // 表单实例，用于添加数据
    const [form] = Form.useForm();
    // 用于地址列搜索输入框的引用
    const searchInput = useRef(null);

    // 自定义地址列的筛选功能，注意：传入的 dataIndex 为 "target"，即后端过滤地址时用的参数名
    const getColumnSearchProps = (dataIndex) => ({
        filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters }) => (
            <div style={{ padding: 8 }} onKeyDown={(e) => e.stopPropagation()}>
                <Input
                    ref={searchInput}
                    placeholder={`搜索 地址`}
                    value={selectedKeys[0]}
                    onChange={e =>
                        setSelectedKeys(e.target.value ? [e.target.value] : [])
                    }
                    onPressEnter={() => {
                        confirm();
                        // 更新 tableParams.filters 中对应的 key
                        setTableParams(prev => ({
                            ...prev,
                            filters: { ...prev.filters, [dataIndex]: selectedKeys }
                        }));
                    }}
                    style={{ marginBottom: 8, display: 'block' }}
                />
                <Button
                    type="primary"
                    onClick={() => {
                        confirm();
                        setTableParams(prev => ({
                            ...prev,
                            filters: { ...prev.filters, [dataIndex]: selectedKeys }
                        }));
                    }}
                    icon={<SearchOutlined />}
                    size="small"
                    style={{ width: 90, marginRight: 8 }}
                >
                    搜索
                </Button>
                <Button
                    onClick={() => {
                        clearFilters();
                        setTableParams(prev => ({
                            ...prev,
                            filters: { ...prev.filters, [dataIndex]: [] }
                        }));
                    }}
                    size="small"
                    style={{ width: 90 }}
                >
                    重置
                </Button>
            </div>
        ),
        filterIcon: filtered => (
            <SearchOutlined style={{ color: filtered ? '#1890ff' : undefined }} />
        ),
        // 这里 onFilter 不在前端过滤，实际过滤在后端处理
        render: text =>
            <Highlighter
                highlightStyle={{ backgroundColor: '#ffc069', padding: 0 }}
                searchWords={
                    (tableParams.filters[dataIndex] && tableParams.filters[dataIndex][0]) ?
                        [tableParams.filters[dataIndex][0]] : []
                }
                autoEscape
                textToHighlight={text ? text.toString() : ''}
            />,
    });

    // 根据 tableParams 请求后端查询数据，并将分页、筛选、排序条件发送给后端
    const fetchData = () => {
        setLoading(true);
        // 将 filters 对象中的数组转换为简单的 key-value 对象（只取第一个值）
        const filterParams = {};
        if (tableParams.filters) {
            Object.keys(tableParams.filters).forEach(key => {
                if (tableParams.filters[key] && tableParams.filters[key].length > 0) {
                    filterParams[key] = tableParams.filters[key][0];
                }
            });
        }
        const params = {
            page: tableParams.pagination.current,
            pageSize: tableParams.pagination.pageSize,
            sortOrder: tableParams.sortOrder,
            sortField: tableParams.sortField,
            ...filterParams,
        };
        fetch(`http://localhost:8080/QueryWBData?${qs.stringify(params)}`)
            .then((res) => res.json())
            .then((resData) => {
                // 假设后端返回的数据中包含 totalCount 或直接返回数组
                setData(resData);
                setTableParams(prev => ({
                    ...prev,
                    pagination: {
                        ...prev.pagination,
                        total: resData.totalCount || resData.length,
                    },
                }));
            })
            .catch(() => {
                message.error("查询数据失败");
            })
            .finally(() => {
                setLoading(false);
            });
    };

    // 添加数据，调用后端接口添加数据
    const handleAdd = (values) => {
        setLoading(true);
        fetch('http://localhost:8080/AddWBData', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            // values 包含 addr, isWhite, valueType
            body: JSON.stringify(values),
        })
            .then((res) => res.json())
            .then(() => {
                message.success("添加成功");
                form.resetFields();
                fetchData();
            })
            .catch(() => {
                message.error("添加失败");
            })
            .finally(() => {
                setLoading(false);
            });
    };

    // 删除数据，调用后端接口删除数据，通过 id 删除记录
    const handleDelete = (id) => {
        setLoading(true);
        fetch(`http://localhost:8080/DeleteWBData?id=${id}`, { method: 'GET' })
            .then((res) => {
                if (res.ok) {
                    message.success("删除成功");
                    fetchData();
                } else {
                    message.error("删除失败");
                }
            })
            .catch(() => {
                message.error("删除失败");
            })
            .finally(() => {
                setLoading(false);
            });
    };

    // 表格列配置，包含展示数据、筛选、排序和操作列（删除按钮）
    const columns = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: '10%',
            sorter: true,
            showSorterTooltip: { target: 'full-header' },
        },
        {
            title: '地址',
            dataIndex: 'addr',
            key: 'addr',
            width: '30%',
            // 调用自定义筛选函数，使用 "target" 作为筛选键（后端使用该参数过滤地址）
            ...getColumnSearchProps("target"),
        },
        {
            title: '是否白名单',
            dataIndex: 'isWhite',
            key: 'isWhite',
            width: '15%',
            render: (isWhite) => (isWhite ? '白名单' : '黑名单'),
            filters: [
                { text: '白名单', value: true },
                { text: '黑名单', value: false },
            ],
        },
        {
            title: '地址类型',
            dataIndex: 'valueType',
            key: 'valueType',
            width: '25%',
            render: (valueType) => {
                switch (valueType) {
                    case 0:
                        return 'BlackIPv4Type';
                    case 1:
                        return 'WhiteIPv4Type';
                    case 2:
                        return 'BlackIPv4NetworkType';
                    case 3:
                        return 'WhiteIPv4NetworkType';
                    case 4:
                        return 'BlackPortType';
                    default:
                        return valueType;
                }
            },
            filters: [
                { text: 'BlackIPv4Type', value: 0 },
                { text: 'WhiteIPv4Type', value: 1 },
                { text: 'BlackIPv4NetworkType', value: 2 },
                { text: 'WhiteIPv4NetworkType', value: 3 },
                { text: 'BlackPortType', value: 4 },
            ],
        },
        {
            title: '操作',
            key: 'action',
            render: (text, record) => (
                <Button type="link" onClick={() => handleDelete(record.id)}>
                    删除
                </Button>
            ),
        },
    ];

    // 表格参数变化时触发查询，将分页、过滤和排序参数更新到 tableParams 中
    const handleTableChange = (pagination, filters, sorter) => {
        setTableParams({
            pagination,
            filters,
            sortOrder: Array.isArray(sorter) ? undefined : sorter.order,
            sortField: Array.isArray(sorter) ? undefined : sorter.field,
        });
        if (pagination.pageSize !== tableParams.pagination.pageSize) {
            setData([]);
        }
    };

    // 当分页、过滤或排序参数发生变化时重新查询数据
    useEffect(() => {
        fetchData();
    }, [
        tableParams.pagination.current,
        tableParams.pagination.pageSize,
        tableParams.sortOrder,
        tableParams.sortField,
        JSON.stringify(tableParams.filters),
    ]);

    return (
        <div>
            {/* 添加数据的表单 */}
            <Form
                form={form}
                layout="inline"
                onFinish={handleAdd}
                style={{ marginBottom: 16 }}
            >
                <Form.Item
                    name="addr"
                    label="地址"
                    rules={[{ required: true, message: '请输入地址' }]}
                >
                    <Input placeholder="请输入地址" />
                </Form.Item>
                <Form.Item
                    name="isWhite"
                    label="白黑名单"
                    valuePropName="checked"
                    initialValue={true}
                >
                    <Switch checkedChildren="白" unCheckedChildren="黑" />
                </Form.Item>
                <Form.Item
                    name="valueType"
                    label="地址类型"
                    rules={[{ required: true, message: '请选择地址类型' }]}
                >
                    <Select style={{ width: 180 }} placeholder="请选择地址类型">
                        <Option value={0}>BlackIPv4Type</Option>
                        <Option value={1}>WhiteIPv4Type</Option>
                        <Option value={2}>BlackIPv4NetworkType</Option>
                        <Option value={3}>WhiteIPv4NetworkType</Option>
                        <Option value={4}>BlackPortType</Option>
                    </Select>
                </Form.Item>
                <Form.Item>
                    <Button type="primary" htmlType="submit">
                        添加
                    </Button>
                </Form.Item>
            </Form>

            {/* 数据展示表格 */}
            <Table
                columns={columns}
                rowKey={(record) => record.id}
                dataSource={data}
                pagination={tableParams.pagination}
                loading={loading}
                onChange={handleTableChange}
                showSorterTooltip={{ target: 'sorter-icon' }}
            />
        </div>
    );
}
