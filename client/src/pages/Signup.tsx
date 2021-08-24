import { Form, Input, Button } from 'antd';
import { Row } from "antd";
import  axios from "axios";

const Signup = () => {
  const onFinish = (values: any) => {
    // post to API here
    const API_URL = `${process.env.REACT_APP_API_ENDPOINT}/signup`;
    axios.post(API_URL, JSON.stringify(values), {
        headers: {
            // 'application/json' is the modern content-type for JSON, but some
            // older servers may use 'text/json'.
            // See: http://bit.ly/text-json
            'Content-Type': 'application/json',
        }
    }).then(response => console.log(response))
    console.log('Success:', values, API_URL);
  };

  const onFinishFailed = (errorInfo: any) => {
    console.log('Failed:', errorInfo);
  };

  return (
    <Row justify="center" align="middle" style={{minHeight: '100vh'}}>
        <Form
        id="login-form"
        name="basic"
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 16 }}
        initialValues={{ remember: true }}
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        >
        <Form.Item
            label="Email"
            name="email"
            rules={[{ required: true, message: 'Please input your username!' }]}
        >
            <Input />
        </Form.Item>

        <Form.Item
            label="Password"
            name="password"
            rules={[{ required: true, message: 'Please input your password!' }]}
        >
            <Input.Password />
        </Form.Item>

        <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
            <Button type="primary" htmlType="submit">
            Submit
            </Button>
        </Form.Item>
        </Form>
    </Row>
  );
};

export default Signup;