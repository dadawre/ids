"use client";
import React, { useEffect, useState } from 'react';
import {Button, Flex, Table} from 'antd';
import qs from 'qs';
import ConfigModelForm from "@/app/component/configData/configModelForm"
const columns = [
    {
        title: '索引',
        dataIndex: 'index',
        sorter: true,
        // render: (name) => `${name.first} ${name.last}`,
        width: '20%',
    },
    {
        title: '配置键',
        dataIndex: 'key',
        width: '20%',
    },
    {
        title: '配置值',
        dataIndex: 'value',
    },
];

const getRandomuserParams = (params) => ({
    results: params.pagination?.pageSize,
    page: params.pagination?.current,
    ...params,
});

export default function ConfigDataTable() {
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);
    const [selectedRowKeys, setSelectedRowKeys] = useState([]);
    const [tableParams, setTableParams] = useState({
        pagination: {
            current: 1,
            pageSize: 10,
        },
    });

    //某行被选择时触发的回调函数
    const onSelectChange = (newSelectedRowKeys) => {
        console.log('selectedRowKeys changed: ', newSelectedRowKeys);
        setSelectedRowKeys(newSelectedRowKeys);
    };

    //选择参数对象
    const rowSelection = {
        selectedRowKeys:selectedRowKeys,
        onChange: onSelectChange,
        type:"radio",
    };

    //是否有选择的行
    const hasSelected = selectedRowKeys.length > 0;


    const fetchData = () => {
        setLoading(true);
        fetch(`http://localhost:8080/configData?${qs.stringify(getRandomuserParams(tableParams))}`)
            .then((res) => res.json())
            .then(({ results, total }) => { // 现在假设后端返回了 total 字段
                setData(results); // Assuming 'results' is an array
                setLoading(false);
                setTableParams((prev) => ({
                    ...prev,
                    pagination: {
                        ...prev.pagination,
                        total: total, // 从响应数据中获取 total
                    },
                }));
            });
    };


    useEffect(() => {
        fetchData();
    }, [
        tableParams.pagination?.current,
        tableParams.pagination?.pageSize,
        tableParams?.sortOrder,
        tableParams?.sortField,
        tableParams.filters, // Ensure filters are passed correctly
    ]);

    const handleTableChange = (pagination, filters, sorter) => {
        setTableParams({
            pagination,
            filters,
            sortOrder: Array.isArray(sorter) ? undefined : sorter.order,
            sortField: Array.isArray(sorter) ? undefined : sorter.field,
        });

        // Reset data if pageSize changes
        if (pagination.pageSize !== tableParams.pagination?.pageSize) {
            setData([]);
        }
    };

    return (
        <Flex gap="middle" vertical>
            <Flex align="center" gap="middle">
                <ConfigModelForm data={data} onChange={setData}></ConfigModelForm>
                {hasSelected ? `Selected ${selectedRowKeys.length} items` : null}
            </Flex>
        <Table
            columns={columns}
            rowKey={(record) => record.index}
            dataSource={data}
            pagination={tableParams.pagination}
            loading={loading}
            onChange={handleTableChange}
            rowSelection={rowSelection}
        />
        </Flex>
    );
}
